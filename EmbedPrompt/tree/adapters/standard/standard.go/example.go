

// factory call here.
func New(basePath string, embedDir string) deps.Deps {
	// ... init other deps
	adapter.Deps.EmbedDeps = EmbedFactory(adapter, embedDir)
	return adapter.Deps
}
