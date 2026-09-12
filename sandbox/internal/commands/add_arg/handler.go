package add_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_arg"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addArgAction.AddArg(sandbox, api.FieldProps{
		Path:        entries.Path,
		Command:     entries.Command,
		Name:        entries.Name,
		Identifiers: nil,
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
