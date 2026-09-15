package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/availableconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AvailablesDir holds one selection per available: which adapter wins for
// each field of Deps. adapters/libs/ is what the project has; this is what it
// binds.
const AvailablesDir = "adapters/availables"

// AvailableConfFile is the declaration of one available. An available
// directory that has one gets its new.go generated from it; one that has none
// is a hand-written mix and is left alone.
const AvailableConfFile = "available.yaml"

// StandardAvailable is the available every project starts with, the one
// cmd/main/main.go imports and the one an install writes into unless the
// caller names another.
const StandardAvailable = "standard"

// AvailableDir is the project-relative directory of one available.
func AvailableDir(available string) string {
	return AvailablesDir + "/" + available
}

// AvailableConfPath is the project-relative path of one available's
// declaration.
func AvailableConfPath(available string) string {
	return AvailableDir(available) + "/" + AvailableConfFile
}

// LoadAvailableConf reads one available's declaration back through the
// transaction-aware io, so a selection written earlier in the same command is
// visible before Persist.
func LoadAvailableConf(sandbox *api.Sandbox, io *smartio.SmartIO, available string) (*availableconf.AvailableConf, error) {
	rel := AvailableConfPath(available)

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not read %s: no available named %q", rel, available)
	}

	return availableconf.New(sandbox, string(content))
}

// DeclaredAvailables returns the name of every available that declares its
// selection, in listing order. An available with no declaration is not
// reported: it is hand-written, and nothing generated may touch it.
func DeclaredAvailables(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var availables []string

	if !io.IsDir(AvailablesDir) {
		return availables
	}

	for _, dir := range io.ListDirs(AvailablesDir) {
		parts := sandbox.Deps.Stringsdeps.Split(dir, "/")
		name := parts[len(parts)-1]
		if name == "" {
			continue
		}
		if _, err := io.ReadFile(AvailableConfPath(name)); err != nil {
			continue
		}
		availables = append(availables, name)
	}

	return availables
}

// EnrollAdapter adds adapter to the selection of every available that declares
// one, because the invariant is that every available fills every field of Deps:
// a contract installed into a project and bound by no available is a nil func
// waiting to panic. A project with no declared available gets the standard
// one, which is what cmd/main/main.go imports.
func EnrollAdapter(sandbox *api.Sandbox, io *smartio.SmartIO, adapter string) error {
	availables := DeclaredAvailables(sandbox, io)
	if len(availables) == 0 {
		return writeAvailable(sandbox, io, StandardAvailable, availableconf.NewEmpty(sandbox), adapter, true)
	}

	for _, name := range availables {
		conf, err := LoadAvailableConf(sandbox, io, name)
		if err != nil {
			return err
		}
		if err := writeAvailable(sandbox, io, name, conf, adapter, true); err != nil {
			return err
		}
	}

	return nil
}

// UnenrollAdapter is EnrollAdapter's inverse: it drops adapter from every
// available that binds it, which is what makes removing an adapter leave a
// tree that still compiles.
func UnenrollAdapter(sandbox *api.Sandbox, io *smartio.SmartIO, adapter string) error {
	for _, name := range DeclaredAvailables(sandbox, io) {
		conf, err := LoadAvailableConf(sandbox, io, name)
		if err != nil {
			return err
		}
		if err := writeAvailable(sandbox, io, name, conf, adapter, false); err != nil {
			return err
		}
	}

	return nil
}

// writeAvailable applies one enrollment change to a selection and writes it
// back only when something actually changed, so a re-install rewrites nothing.
func writeAvailable(sandbox *api.Sandbox, io *smartio.SmartIO, name string, conf *availableconf.AvailableConf, adapter string, enroll bool) error {
	changed := conf.Remove(adapter)
	if enroll {
		changed = conf.Add(adapter)
	}

	if !changed {
		return nil
	}

	return io.WriteFileOverwrite(AvailableConfPath(name), []byte(conf.Render()))
}

// SelectAdapter makes adapter the one this available binds for the dep it
// fills, dropping whichever other adapter filled that same field. It is the
// single place the "exactly one adapter per field per available" invariant is
// maintained, so a contract with two implementations can never end up with
// both bound.
func SelectAdapter(sandbox *api.Sandbox, io *smartio.SmartIO, available string, adapter string) error {

	conf, err := LoadAvailableConf(sandbox, io, available)
	if err != nil {
		return err
	}

	target, err := LoadAdapterConf(sandbox, io, adapter)
	if err != nil {
		return err
	}

	bound := append([]string{}, conf.Adapters...)
	changed := conf.Add(adapter)

	for _, other := range bound {
		if other == adapter {
			continue
		}
		other_conf, err := LoadAdapterConf(sandbox, io, other)
		if err != nil || other_conf.Dep != target.Dep {
			continue
		}
		if conf.Remove(other) {
			changed = true
		}
	}

	if !changed {
		return nil
	}

	return io.WriteFileOverwrite(AvailableConfPath(available), []byte(conf.Render()))
}

// AvailablesBinding returns the name of every available that binds adapter, in
// listing order. It is what `remove-adapter` refuses on and what
// `list-adapters` shows.
func AvailablesBinding(sandbox *api.Sandbox, io *smartio.SmartIO, adapter string) []string {
	var binding []string

	for _, name := range DeclaredAvailables(sandbox, io) {
		conf, err := LoadAvailableConf(sandbox, io, name)
		if err != nil {
			continue
		}
		if conf.Has(adapter) {
			binding = append(binding, name)
		}
	}

	return binding
}

// ValidateAvailableName rejects a name that could not be a directory under
// adapters/availables/ and a Go package clause at the same time — the same
// check ValidateCommandName makes, for the same reason: the name is
// propagated straight into `package <name>` of the generated new.go.
func ValidateAvailableName(sandbox *api.Sandbox, available string) error {
	if available == "" {
		return sandbox.Deps.Std.Errorf("an available needs a name")
	}
	if available[0] < 'a' || available[0] > 'z' {
		return sandbox.Deps.Std.Errorf("invalid available name %q: an available name must start with a lowercase letter", available)
	}
	for _, letter := range available {
		if (letter >= 'a' && letter <= 'z') || (letter >= '0' && letter <= '9') {
			continue
		}
		return sandbox.Deps.Std.Errorf("invalid available name %q: only lowercase letters and digits are allowed (it becomes the directory %s and a Go package name)",
			available, AvailableDir(available))
	}
	return nil
}
