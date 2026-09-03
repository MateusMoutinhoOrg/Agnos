package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/set_command"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	set_error := setCommandAction.SetCommand(deps, api.CommandProps{
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
		deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
