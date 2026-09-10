package moduleconf

type ModuleConf struct {
	Module    string
	GoVersion string
	Requires  []string

	// Directives are the lines of every other go.mod directive — `replace`,
	// `exclude`, `toolchain`, `retract` — kept verbatim and re-emitted
	// unchanged. Nothing here reads them, and dropping what it does not read
	// would silently unwire a project developing two repos side by side.
	Directives []string

	// AddRequire inserts a "<module> <version>" require entry, replacing any
	// existing entry for the same module path.
	AddRequire func(require string)
	// RemoveRequire drops every require entry for the given module path.
	RemoveRequire func(module string)

	Render func() string
}
