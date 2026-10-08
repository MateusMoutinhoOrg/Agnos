package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/bindingconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// BindingsDir holds one selection per binding: which adapter wins for
// each field of Deps. adapters/impls/ is what the project has; this is what it
// binds.
const BindingsDir = "adapters/bindings"

// BindingConfFile is the declaration of one binding. A binding
// directory that has one gets its new.go generated from it; one that has none
// is a hand-written mix and is left alone.
const BindingConfFile = "binding.yaml"

// StandardBinding is the binding every project starts with, the one
// cmd/main/main.go imports and the one an install writes into unless the
// caller names another.
const StandardBinding = "standard"

// BindingDir is the project-relative directory of one binding.
func BindingDir(binding string) string {
	return BindingsDir + "/" + binding
}

// LegacyBindingsDir is where a repo rendered by an agnos older than the
// binding rename keeps its bindings, under their old name of availables. A
// remote dep may be such a repo, and its tree is not ours to rename.
const LegacyBindingsDir = "adapters/availables"

// RemoteBindingDir is the module-relative directory of one binding of a remote
// repo checked out at moduleDir — adapters/bindings/<binding>, or the legacy
// adapters/availables/<binding> — and whether the repo has it at all.
func RemoteBindingDir(sandbox *api.Sandbox, moduleDir string, binding string) (string, bool) {
	for _, dir := range []string{BindingDir(binding), LegacyBindingsDir + "/" + binding} {
		if sandbox.Deps.IoDeps.IsDir(sandbox.Deps.IoDeps.Join(moduleDir, dir)) {
			return dir, true
		}
	}
	return BindingDir(binding), false
}

// BindingConfPath is the project-relative path of one binding's
// declaration.
func BindingConfPath(binding string) string {
	return BindingDir(binding) + "/" + BindingConfFile
}

// LoadBindingConf reads one binding's declaration back through the
// transaction-aware io, so a selection written earlier in the same command is
// visible before Persist.
func LoadBindingConf(sandbox *api.Sandbox, io *stagedfs.StagedFS, binding string) (*bindingconf.BindingConf, error) {
	rel := BindingConfPath(binding)

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("could not read %s: no binding named %q", rel, binding)
	}

	return bindingconf.New(sandbox, string(content))
}

// DeclaredBindings returns the name of every binding that declares its
// selection, in listing order. A binding with no declaration is not
// reported: it is hand-written, and nothing generated may touch it.
func DeclaredBindings(sandbox *api.Sandbox, io *stagedfs.StagedFS) []string {
	var bindings []string

	if !io.IsDir(BindingsDir) {
		return bindings
	}

	for _, dir := range io.ListDirs(BindingsDir) {
		parts := sandbox.Deps.StringsDeps.Split(dir, "/")
		name := parts[len(parts)-1]
		if name == "" {
			continue
		}
		if _, err := io.ReadFile(BindingConfPath(name)); err != nil {
			continue
		}
		bindings = append(bindings, name)
	}

	return bindings
}

// EnrollAdapter adds adapter to the selection of every binding that declares
// one, because the invariant is that every binding fills every field of Deps:
// a contract installed into a project and bound by no binding is a nil func
// waiting to panic. A project with no declared binding gets the standard
// one, which is what cmd/main/main.go imports.
func EnrollAdapter(sandbox *api.Sandbox, io *stagedfs.StagedFS, adapter string) error {
	bindings := DeclaredBindings(sandbox, io)
	if len(bindings) == 0 {
		return writeBinding(sandbox, io, StandardBinding, bindingconf.NewEmpty(sandbox), adapter, true)
	}

	for _, name := range bindings {
		conf, err := LoadBindingConf(sandbox, io, name)
		if err != nil {
			return err
		}
		if err := writeBinding(sandbox, io, name, conf, adapter, true); err != nil {
			return err
		}
	}

	return nil
}

// UnenrollAdapter is EnrollAdapter's inverse: it drops adapter from every
// binding that binds it, which is what makes removing an adapter leave a
// tree that still compiles.
func UnenrollAdapter(sandbox *api.Sandbox, io *stagedfs.StagedFS, adapter string) error {
	for _, name := range DeclaredBindings(sandbox, io) {
		conf, err := LoadBindingConf(sandbox, io, name)
		if err != nil {
			return err
		}
		if err := writeBinding(sandbox, io, name, conf, adapter, false); err != nil {
			return err
		}
	}

	return nil
}

// writeBinding applies one enrollment change to a selection and writes it
// back only when something actually changed, so a re-install rewrites nothing.
func writeBinding(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, conf *bindingconf.BindingConf, adapter string, enroll bool) error {
	changed := conf.Remove(adapter)
	if enroll {
		changed = conf.Add(adapter)
	}

	if !changed {
		return nil
	}

	return io.WriteFile(BindingConfPath(name), []byte(conf.Render()))
}

// SelectAdapter makes adapter the one this binding binds for the dep it
// fills, dropping whichever other adapter filled that same field. It is the
// single place the "exactly one adapter per field per binding" invariant is
// maintained, so a contract with two implementations can never end up with
// both bound.
func SelectAdapter(sandbox *api.Sandbox, io *stagedfs.StagedFS, binding string, adapter string) error {

	conf, err := LoadBindingConf(sandbox, io, binding)
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

	return io.WriteFile(BindingConfPath(binding), []byte(conf.Render()))
}

// BindingsUsing returns the name of every binding that binds adapter, in
// listing order. It is what `remove-adapter` refuses on and what
// `list-adapters` shows.
func BindingsUsing(sandbox *api.Sandbox, io *stagedfs.StagedFS, adapter string) []string {
	var binding []string

	for _, name := range DeclaredBindings(sandbox, io) {
		conf, err := LoadBindingConf(sandbox, io, name)
		if err != nil {
			continue
		}
		if conf.Has(adapter) {
			binding = append(binding, name)
		}
	}

	return binding
}

// ValidateBindingName rejects a name that could not be a directory under
// adapters/bindings/ and a Go package clause at the same time — the same
// check ValidateCommandName makes, for the same reason: the name is
// propagated straight into `package <name>` of the generated new.go.
func ValidateBindingName(sandbox *api.Sandbox, binding string) error {
	if binding == "" {
		return sandbox.Deps.StdDeps.Errorf("a binding needs a name")
	}
	if binding[0] < 'a' || binding[0] > 'z' {
		return sandbox.Deps.StdDeps.Errorf("invalid binding name %q: a binding name must start with a lowercase letter", binding)
	}
	for _, letter := range binding {
		if (letter >= 'a' && letter <= 'z') || (letter >= '0' && letter <= '9') {
			continue
		}
		return sandbox.Deps.StdDeps.Errorf("invalid binding name %q: only lowercase letters and digits are allowed (it becomes the directory %s and a Go package name)",
			binding, BindingDir(binding))
	}
	return nil
}
