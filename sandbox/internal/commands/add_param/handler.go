package add_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_param"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addParamAction.AddParam(sandbox, api.RouteFieldProps{
		Path:        entries.Path,
		Route:       entries.Route,
		Name:        entries.Name,
		Description: entries.Description,
		Examples:    entries.Example,
		Type:        entries.Type,
		Default:     entries.Default,
		Required:    entries.Required,
		Array:       entries.Array,
		Min:         entries.Min,
		Max:         entries.Max,
		Position:    entries.Position,
	})
	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
