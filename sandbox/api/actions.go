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
	Start         func(props StartProps) error
	Build         func(props BuildProps) error
	Install       func(props InstallProps) error
	Uninstall     func(props UninstallProps) error
	List          func(props ListProps) error
	ExtensionHelp func(props ExtensionHelpProps) error
}
