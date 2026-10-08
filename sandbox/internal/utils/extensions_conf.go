package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/extensionsconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// ExtensionsConfFile is the declaration saying which generation mechanics this
// project wants. It is the whole of what `build` reads to decide what to
// render: nothing is inferred from the tree on disk.
const ExtensionsConfFile = "extensions.yaml"

// The extension names, each the name of the mechanic it turns on and of the
// <x>-init / <x>-purge pair that owns it. Every one but doc and readme renders
// into the sandbox and needs it on; those two stand on their own because
// neither needs the sandbox to be generated.
const (
	ExtensionSandbox    = "sandbox"
	ExtensionDeps       = "deps"
	ExtensionCli        = "cli"
	ExtensionServer     = "server"
	ExtensionFront      = "front"
	ExtensionDatabase   = "database"
	ExtensionExample    = "example"
	ExtensionBackoffice = "backoffice"
	ExtensionDoc        = "doc"
	ExtensionReadme     = "readme"
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
		{ExtensionSandbox, true, "the sandbox core: sandbox/new.go, api/sandbox.go, internal/generated/config"},
		{ExtensionDeps, false, "the dependency layer: sandbox/deps/, adapters/, bindings"},
		{ExtensionCli, false, "the cli layer: cmd/main, help, version and the OpinionatedAgnosCli lib"},
		{ExtensionServer, false, "the http layer: server/, routes/ and the OpinionatedAgnosServer lib"},
		{ExtensionFront, false, "the front layer: the route serving assets/front/ and the OpinionatedAgnosFront lib"},
		{ExtensionDatabase, false, "the database layer: the declared databases and the OpinionatedAgnosDatabase lib"},
		{ExtensionExample, true, "the examples/ suite and run-examples"},
		{ExtensionBackoffice, false, "the admin backoffice: login, users, API tokens and the /api/admin JSON api"},
		{ExtensionDoc, true, "the docs/ tree and its Index.md files"},
		{ExtensionReadme, true, "README.md, built from themes.yaml and the doc index"},
	}
}

// ExtensionRequires is what one mechanic needs on to have anything to render
// into: every mechanic but doc and readme renders into the sandbox, the server layer is opened
// by a command of the cli layer (start-server), and the front layer is served
// by a route of the server layer.
func ExtensionRequires(name string) []string {
	switch name {
	case ExtensionServer:
		return []string{ExtensionSandbox, ExtensionCli}
	case ExtensionFront:
		return []string{ExtensionSandbox, ExtensionServer}
	case ExtensionBackoffice:
		return []string{ExtensionSandbox, ExtensionServer, ExtensionFront, ExtensionDatabase}
	case ExtensionSandbox, ExtensionDoc, ExtensionReadme:
		return nil
	}
	return []string{ExtensionSandbox}
}

// ExtensionDependents is every mechanic that requires name.
func ExtensionDependents(name string) []string {
	dependents := []string{}
	for _, spec := range ExtensionCatalog() {
		if contains(ExtensionRequires(spec.Name), name) {
			dependents = append(dependents, spec.Name)
		}
	}
	return dependents
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
	return ConfigDir(sandbox) + "/" + ExtensionsConfFile
}

// LoadExtensionsConf reads AgnosConfig/extensions.yaml back through the
// transaction-aware io. Like project.yaml it is written once by `agnos start`,
// so a missing or unparsable file is a hard error: what a project generates is
// declared, never guessed from the directories it happens to carry.
func LoadExtensionsConf(sandbox *api.Sandbox, io *stagedfs.StagedFS) (*extensionsconf.ExtensionsConf, error) {
	rel := ExtensionsConfPath(sandbox)

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("could not read %s: run `agnos start` first (%w)", rel, err)
	}

	conf, err := extensionsconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("%s: %w", rel, err)
	}
	return conf, nil
}

// SaveExtensionsConf renders conf back over the declaration.
func SaveExtensionsConf(sandbox *api.Sandbox, io *stagedfs.StagedFS, conf *extensionsconf.ExtensionsConf) error {
	return io.WriteFile(ExtensionsConfPath(sandbox), []byte(conf.Render()))
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
func ExtensionEnabled(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) (bool, error) {
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
func SetExtension(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, enabled bool) error {
	if !IsExtensionName(name) {
		return sandbox.Deps.StdDeps.Errorf("unknown extension %q (known: %s)",
			name, sandbox.Deps.StringsDeps.Join(ExtensionNames(), ", "))
	}

	conf, err := LoadExtensionsConf(sandbox, io)
	if err != nil {
		return err
	}

	NormalizeExtensions(conf)

	// A declaration never holds a mechanic on without what it renders into:
	// turning one off while a dependent is on, or one on while a requirement
	// is off, is refused, and verify holds every declaration to the same.
	if enabled {
		for _, required := range ExtensionRequires(name) {
			if !conf.IsEnabled(required) {
				return sandbox.Deps.StdDeps.Errorf("cannot turn %s on: it needs %s, which is off", name, required)
			}
		}
	} else {
		for _, dependent := range ExtensionDependents(name) {
			if conf.IsEnabled(dependent) {
				return sandbox.Deps.StdDeps.Errorf("cannot turn %s off: %s is on and needs it (turn %s off first)", name, dependent, dependent)
			}
		}
	}
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
func RequireExtension(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) error {
	enabled, err := ExtensionEnabled(sandbox, io, name)
	if err != nil {
		return err
	}
	if !enabled {
		return sandbox.Deps.StdDeps.Errorf("the %s extension is off: turn it on with `agnos enable-extension %s`", name, name)
	}
	return nil
}
