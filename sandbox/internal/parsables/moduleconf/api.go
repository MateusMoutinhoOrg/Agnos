package moduleconf

type ModuleConf struct {
	Module    string
	GoVersion string
	Requires  []string

	// AddRequire inserts a "<module> <version>" require entry, replacing any
	// existing entry for the same module path.
	AddRequire func(require string)
	// RemoveRequire drops every require entry for the given module path.
	RemoveRequire func(module string)

	Render func() string
}
