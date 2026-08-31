package api

type StartProps struct {
	Path        string
	ProjectName string
	Module      *string
	Force       bool
}

type BuildProps struct {
	Path    string
	Project string
	Force   bool
}

type InstallProps struct {
	Path string
	Item string
}

type UninstallProps struct {
	Path string
	Item string
}

type ListProps struct {
	Path string
}

type ExtensionHelpProps struct {
	Path      string
	Extension string
}

type Actions struct {
	Build                   func(props BuildProps) error
	InstallExtension        func(props InstallProps) error
	RemoveExtension         func(props UninstallProps) error
	ListAvaliableExtensions func(props ListProps) ([]string, error)
}
