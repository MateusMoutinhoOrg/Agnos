package adapterconf

// OriginCatalog marks an adapter that came out of the embedded catalog,
// assets/adapterlist/<adapter>/.
const OriginCatalog = "catalog"

// OriginGenerated marks an adapter written by the generator itself — the shim
// of a dep copied from a remote agnos repo. `remove-adapter` refuses one:
// half of that dep is the contract beside it, so `remove-dep` owns the pair.
const OriginGenerated = "generated"

// AdapterConf is the parsed form of one adapter declaration: the adapter.yaml
// at the root of assets/adapterlist/<adapter>/, and the copy that install
// writes to adapters/libs/<adapter>/adapter.yaml. That copy is what tells
// which dep an adapter fills without parsing the body of its Bind.
type AdapterConf struct {
	// Name is the adapter, the same spelling as its catalog directory and as
	// the adapters/libs/<adapter>/ it installs.
	Name string

	// Dep is the dep whose contract this adapter implements — the directory
	// under sandbox/deps/, so the Deps field is its title-cased spelling.
	Dep string

	// Help is the one-line description `list-adapters` prints.
	Help string

	// Module is the versioned module this adapter imports, written as
	// "<module>@<version>", or "" when it needs none beyond the stdlib.
	Module string

	// Origin is OriginCatalog or OriginGenerated.
	Origin string

	// ModuleSpec splits Module into its path and version, and reports whether
	// the adapter pins one at all.
	ModuleSpec func() (module string, version string, ok bool)

	Render func() string
}
