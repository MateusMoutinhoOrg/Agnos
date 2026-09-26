package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/extensionsconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ExtensionsConfFile is the declaration saying which generation mechanics this
// project wants. It is the whole of what `build` reads to decide what to
// render: nothing is inferred from the tree on disk.
const ExtensionsConfFile = "extensions.yaml"

// The extension names. Every mechanic that depends on the sandbox is spelled
// sandbox-<mechanic>; doc and readme stand on their own because neither needs
// the sandbox to be generated.
const (
	ExtensionSandbox         = "sandbox"
	ExtensionSandboxDeps     = "sandbox-deps"
	ExtensionSandboxCli      = "sandbox-cli"
	ExtensionSandboxServer   = "sandbox-server"
	ExtensionSandboxFront    = "sandbox-front"
	ExtensionSandboxDatabase = "sandbox-database"
	ExtensionSandboxExample  = "sandbox-example"
	ExtensionDoc             = "doc"
	ExtensionReadme          = "readme"
)

// ExtensionSpec is one mechanic of the catalog: its key, what a fresh project
// starts with, and the one line `list-extensions` prints for it.
type ExtensionSpec struct {
	Name    string
	Default bool
	Help    string
}

// ExtensionCatalog is every mechanic agnos knows, in the order they are listed
// and rendered. A key outside this list is a violation (see CheckExtensions),
// and a key of this list missing from a project's declaration is filled with
// its Default by NormalizeExtensions.
func ExtensionCatalog() []ExtensionSpec {
	return []ExtensionSpec{
		{ExtensionSandbox, true, "the sandbox core: sandbox/new.go, api/sandbox.go, internal/config"},
		{ExtensionSandboxDeps, false, "the dependency layer: sandbox/deps/, adapters/, availables"},
		{ExtensionSandboxCli, false, "the cli layer: cmd/main, the dispatch, help and version"},
		{ExtensionSandboxServer, false, "the http layer: server/, routes/, routeio/"},
		{ExtensionSandboxFront, false, "the front layer: frontio/ and the route serving assets/frontend/"},
		{ExtensionSandboxDatabase, false, "the database layer: databaseio/ and the declared databases"},
		{ExtensionSandboxExample, true, "the examples/ suite and exec-test"},
		{ExtensionDoc, true, "the docs/ tree and its Index.md files"},
		{ExtensionReadme, true, "README.md, built from themes.yaml and the doc index"},
	}
}

// ExtensionNames is ExtensionCatalog reduced to its keys.
func ExtensionNames() []string {
	specs := ExtensionCatalog()
	names := make([]string, 0, len(specs))
	for _, spec := range specs {
		names = append(names, spec.Name)
	}
	return names
}

// IsExtensionName reports whether name is a key of the catalog.
func IsExtensionName(name string) bool {
	for _, spec := range ExtensionCatalog() {
		if spec.Name == name {
			return true
		}
	}
	return false
}

// ExtensionsConfPath is the project-relative path of the extensions
// declaration, held in the config directory beside project.yaml.
func ExtensionsConfPath(sandbox *api.Sandbox) string {
	return sandbox.Config.ProjectName + "Config/" + ExtensionsConfFile
}

// LoadExtensionsConf reads <ProjectName>Config/extensions.yaml back through the
// transaction-aware io. Like project.yaml it is written once by `agnos start`,
// so a missing or unparsable file is a hard error: what a project generates is
// declared, never guessed from the directories it happens to carry.
func LoadExtensionsConf(sandbox *api.Sandbox, io *smartio.SmartIO) (*extensionsconf.ExtensionsConf, error) {
	rel := ExtensionsConfPath(sandbox)

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not read %s: run `agnos start` first (%w)", rel, err)
	}

	conf, err := extensionsconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("%s: %w", rel, err)
	}
	return conf, nil
}

// SaveExtensionsConf renders conf back over the declaration.
func SaveExtensionsConf(sandbox *api.Sandbox, io *smartio.SmartIO, conf *extensionsconf.ExtensionsConf) error {
	return io.WriteFileOverwrite(ExtensionsConfPath(sandbox), []byte(conf.Render()))
}

// NewExtensionsConf is the declaration a fresh project starts with: every key
// of the catalog at its default.
func NewExtensionsConf(sandbox *api.Sandbox) *extensionsconf.ExtensionsConf {
	conf := extensionsconf.NewEmpty(sandbox)
	for _, spec := range ExtensionCatalog() {
		conf.SetEnabled(spec.Name, spec.Default)
	}
	return conf
}

// NormalizeExtensions fills every catalog key the declaration does not carry
// with that key's default, and reports whether anything was added. `build`
// writes the file back when it was, so a project scaffolded before a mechanic
// existed gains the key instead of failing on it.
func NormalizeExtensions(conf *extensionsconf.ExtensionsConf) bool {
	changed := false
	for _, spec := range ExtensionCatalog() {
		if !conf.Has(spec.Name) {
			conf.SetEnabled(spec.Name, spec.Default)
			changed = true
		}
	}
	return changed
}

// ExtensionEnabled reads one key of a project's declaration. It is what an
// action asks when it needs to know whether a mechanic is on — never the
// directories that mechanic happens to have written.
func ExtensionEnabled(sandbox *api.Sandbox, io *smartio.SmartIO, name string) (bool, error) {
	conf, err := LoadExtensionsConf(sandbox, io)
	if err != nil {
		return false, err
	}
	NormalizeExtensions(conf)
	return conf.IsEnabled(name), nil
}

// SetExtension flips one key of a project's declaration and writes it back. It
// is what `enable-extension`, `disable-extension` and every <x>-init / <x>-purge
// pair go through, so the declaration is never edited by hand.
func SetExtension(sandbox *api.Sandbox, io *smartio.SmartIO, name string, enabled bool) error {
	if !IsExtensionName(name) {
		return sandbox.Deps.Std.Errorf("unknown extension %q (known: %s)",
			name, sandbox.Deps.Stringsdeps.Join(ExtensionNames(), ", "))
	}

	conf, err := LoadExtensionsConf(sandbox, io)
	if err != nil {
		return err
	}

	NormalizeExtensions(conf)
	conf.SetEnabled(name, enabled)

	if err := SaveExtensionsConf(sandbox, io, conf); err != nil {
		return err
	}

	// Turning a mechanic off removes nothing: agnos simply stops rendering
	// what it owns. Turning one on renders its code group here, so the build
	// that follows collects against a tree that already has the contracts the
	// mechanic contributes to sandbox/api/.
	if !enabled {
		return nil
	}

	return RenderExtensionCode(sandbox, io, name)
}

// RequireExtension fails when one mechanic is off, with the command that turns
// it back on. It is what a command owned by a mechanic asks before it does
// anything: a project that declared the mechanic off is not asking agnos to
// look after those files.
func RequireExtension(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	enabled, err := ExtensionEnabled(sandbox, io, name)
	if err != nil {
		return err
	}
	if !enabled {
		return sandbox.Deps.Std.Errorf("the %s extension is off: turn it on with `agnos enable-extension %s`", name, name)
	}
	return nil
}
