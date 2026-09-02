package depsversionconf

// DepsVersionConf is the parsed form of assets/depsversion.yaml: a mapping
// from a dep name (the same name used by dep-install / dep-remove) to the
// versioned module path that dep pulls, written as "<module>@<version>".
type DepsVersionConf struct {
	// Deps maps dep name -> "<module>@<version>".
	Deps map[string]string

	// Get returns the module path and version pinned for dep, and whether
	// dep is listed at all.
	Get func(dep string) (module string, version string, ok bool)

	Render func() string
}
