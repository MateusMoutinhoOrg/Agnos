package depconf

// DepConf is the parsed form of one dep declaration: the dep.yaml at the root
// of assets/deplist/<dep>/, beside the sandbox/deps/<dep>/ tree that dep
// installs. It names the contract, the Deps field that contract fills, and the
// adapter `add-dep` installs when the caller chooses none.
type DepConf struct {
	// Name is the dep, the same spelling as its catalog directory and as the
	// sandbox/deps/<dep>/ it installs.
	Name string

	// Field is the Deps field the contract fills, the title-cased Name that
	// sandbox/deps/deps.go is generated with.
	Field string

	// Help is the one-line description `list-deps` prints.
	Help string

	// DefaultAdapter is the adapter installed with the contract when the
	// caller passes no --adapter.
	DefaultAdapter string

	Render func() string
}
