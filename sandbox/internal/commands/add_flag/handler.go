package add_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_flag"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addFlagAction.AddFlag(sandbox, api.FieldProps{
		Path:        entries.Path,
		Command:     entries.Command,
		Name:        entries.Name,
		Identifiers: entries.Identifier,
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
