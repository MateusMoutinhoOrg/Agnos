package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_command"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	set_error := setCommandAction.SetCommand(sandbox, api.CommandProps{
		Path:            entries.Path,
		Command:         entries.Name,
		Help:            entries.Help,
		Category:        entries.Category,
		LongDescription: entries.LongDescription,
		Hidden:          entries.Hidden,
		Visible:         entries.Visible,
		Identifiers:     entries.Identifier,
		Examples:        entries.Example,
	})
	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
