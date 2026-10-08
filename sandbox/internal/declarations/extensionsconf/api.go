package extensionsconf

type Extension struct {
	Name    string
	Enabled bool
}

type ExtensionsConf struct {
	Extensions []Extension

	SetEnabled func(name string, enabled bool)
	IsEnabled  func(name string) bool
	Has        func(name string) bool
	Render     func() string
}
