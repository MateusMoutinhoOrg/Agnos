package api

const (
	RequirementTypeString = iota
	RequirementTypeBytes  = iota
	RequirementTypeInt    = iota
	RequirementTypeFloat  = iota
	RequirementTypeBool   = iota
)

type ExtensionRequirement struct {
	Id       string
	Type     int
	Default  any
	Required bool
}

type ExtensionAction struct {
	Name               string
	EntrieRequirements []ExtensionRequirement
	Handler            func(sandbox *Sandbox, entries map[string]any) error
}

type Extension struct {
	Name     string
	RunBuild bool
	Build    func(sandbox *Sandbox) error
	Actions  []ExtensionAction
}
