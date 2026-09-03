package add_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/add_arg"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	add_error := addArgAction.AddArg(deps, api.FieldProps{
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
		deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
