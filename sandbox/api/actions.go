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
	Build func(props BuildProps) error
}
