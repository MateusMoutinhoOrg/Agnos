package availableconf

// AvailableConf is the parsed form of one available declaration,
// adapters/availables/<name>/available.yaml: the adapters that available
// binds, in the order it binds them. It is a selection, not an inventory —
// adapters/libs/ is what the project has, and this is which of them wins for
// each field of Deps.
type AvailableConf struct {
	// Adapters names one adapters/libs package per entry, in bind order.
	Adapters []string

	// Has reports whether the available already binds adapter.
	Has func(adapter string) bool

	// Add appends adapter, keeping the list sorted and free of duplicates,
	// and reports whether anything changed.
	Add func(adapter string) bool

	// Remove drops adapter and reports whether anything changed.
	Remove func(adapter string) bool

	Render func() string
}
