package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AdapterlistGroup is the embedded catalog of installable adapters: one
// directory per adapter, holding the adapter.yaml that declares it beside the
// adapters/libs/<adapter>/ tree it installs.
const AdapterlistGroup = "adapterlist"

// AdapterConfFile is the declaration of one adapter. It sits at the root of
// the catalog directory, and install writes it into the installed package —
// that copy is what tells which dep an adapter fills without parsing its Bind.
const AdapterConfFile = "adapter.yaml"

// AdaptersDir holds one adapter package per implementation, the open side of
// the pair a contract declares.
const AdaptersDir = "adapters/libs"

// AdapterDir is the project-relative directory of one installed adapter.
func AdapterDir(adapter string) string {
	return AdaptersDir + "/" + adapter
}

// AdapterConfPath is the project-relative path of one installed adapter's
// declaration.
func AdapterConfPath(adapter string) string {
	return AdapterDir(adapter) + "/" + AdapterConfFile
}

// LoadCatalogAdapterConf reads assets/adapterlist/<adapter>/adapter.yaml out of
// the embedded catalog. An adapter with no declaration is not an adapter, so
// the error is the one callers report for an unknown name.
func LoadCatalogAdapterConf(deps *deps.Deps, adapter string) (*adapterconf.AdapterConf, error) {
	content, err := deps.Embeddeps.ReadFile(AdapterlistGroup + "/" + adapter + "/" + AdapterConfFile)
	if err != nil {
		return nil, deps.Std.Errorf("unknown adapter %q", adapter)
	}

	return adapterconf.New(deps, string(content))
}

// LoadAdapterConf reads the declaration of one adapter already installed in
// the target project, through the transaction-aware io.
func LoadAdapterConf(deps *deps.Deps, io *smartio.SmartIO, adapter string) (*adapterconf.AdapterConf, error) {
	rel := AdapterConfPath(adapter)

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, deps.Std.Errorf("could not read %s: %w", rel, err)
	}

	return adapterconf.New(deps, string(content))
}

// CatalogAdapters returns the name of every adapter in the embedded catalog,
// in listing order.
func CatalogAdapters(deps *deps.Deps) ([]string, error) {
	return catalogEntries(deps, AdapterlistGroup)
}

// InstalledAdapters returns the name of every adapter installed in the target
// project, one per adapters/libs sub-directory, in listing order.
func InstalledAdapters(deps *deps.Deps, io *smartio.SmartIO) []string {
	var adapters []string

	if !io.IsDir(AdaptersDir) {
		return adapters
	}

	for _, dir := range io.ListDirs(AdaptersDir) {
		parts := deps.Stringsdeps.Split(dir, "/")
		name := parts[len(parts)-1]
		if name == "" {
			continue
		}
		adapters = append(adapters, name)
	}

	return adapters
}

// AdaptersFillingDep returns the name of every installed adapter whose
// declaration names dep, in listing order. It is the question `remove-dep`
// asks before refusing, and the one an available resolves a field with.
func AdaptersFillingDep(deps *deps.Deps, io *smartio.SmartIO, dep string) []string {
	var adapters []string

	for _, adapter := range InstalledAdapters(deps, io) {
		conf, err := LoadAdapterConf(deps, io, adapter)
		if err != nil {
			continue
		}
		if conf.Dep == dep {
			adapters = append(adapters, adapter)
		}
	}

	return adapters
}
