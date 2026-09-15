package api

type ExtensionAction struct {
	Name    string
	Handler func(sandbox *Sandbox, entries any) error
}

type Extension struct {
	Name string

	Build   func(sandbox *Sandbox) error
	Actions []ExtensionAction
}
