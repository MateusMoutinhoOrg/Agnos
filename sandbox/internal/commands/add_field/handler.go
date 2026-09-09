package add_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_field"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	add_error := addFieldAction.AddField(deps, api.RouteFieldProps{
		Path:        entries.Path,
		Route:       entries.Route,
		Name:        entries.Name,
		Identifier:  entries.Identifier,
		In:          entries.In,
		Description: entries.Description,
		Examples:    entries.Example,
		Type:        entries.Type,
		Default:     entries.Default,
		Required:    entries.Required,
		Array:       entries.Array,
		Min:         entries.Min,
		Max:         entries.Max,
		Position:    entries.Position,
		Format:      entries.Format,
		Pattern:     entries.Pattern,
	})
	if add_error != nil {
		deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
