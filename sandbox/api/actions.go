package api

type BuildProps struct {
	Path  string
	Force bool
}

type InstallProps struct {
	Path string
	Item string
}

type UninstallProps struct {
	Path string
	Item string
}

type ListAvailableExtensionsProps struct {
	Path string
}

type Actions struct {
	Build                   func(props BuildProps) error
	InstallExtension        func(props InstallProps) error
	RemoveExtension         func(props UninstallProps) error
	ListAvaliableExtensions func(props ListAvailableExtensionsProps) ([]string, error)
}
