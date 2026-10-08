package bindingconf

// BindingConf is the parsed form of one binding declaration,
// adapters/bindings/<name>/binding.yaml: the adapters that binding
// binds, in the order it binds them. It is a selection, not an inventory —
// adapters/impls/ is what the project has, and this is which of them wins for
// each field of Deps.
type BindingConf struct {
	// Adapters names one adapters/impls package per entry, in bind order.
	Adapters []string

	// Has reports whether the binding already binds adapter.
	Has func(adapter string) bool

	// Add appends adapter, keeping the list sorted and free of duplicates,
	// and reports whether anything changed.
	Add func(adapter string) bool

	// Remove drops adapter and reports whether anything changed.
	Remove func(adapter string) bool

	Render func() string
}
