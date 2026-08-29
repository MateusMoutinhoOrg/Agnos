package lib

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

// InstallProps carries the operands of the install command: the project
// directory the extension is installed into, and the item to install.
type InstallProps struct {
	Path string
	Item string
}

// UninstallProps carries the operands of the uninstall command: the project
// directory the extension is removed from, and the item to uninstall.
type UninstallProps struct {
	Path string
	Item string
}

// ListProps carries the operands of the list command: the project directory
// whose available extensions are listed.
type ListProps struct {
	Path string
}

// ExtensionHelpProps carries the operands of the extension-help command: the
// project directory and the extension whose help is shown.
type ExtensionHelpProps struct {
	Path      string
	Extension string
}

type CoreApi struct {
	Start         func(props StartProps) error
	Build         func(props BuildProps) error
	Install       func(props InstallProps) error
	Uninstall     func(props UninstallProps) error
	List          func(props ListProps) error
	ExtensionHelp func(props ExtensionHelpProps) error
}
