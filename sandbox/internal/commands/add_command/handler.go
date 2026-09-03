package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/add_command"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	add_error := addCommandAction.AddCommand(deps, entries.Path, entries.Name, entries.Help, entries.Category)

	if add_error != nil {
		deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
