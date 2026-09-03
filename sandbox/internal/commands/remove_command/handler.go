package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_command"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeCommandAction.RemoveCommand(deps, entries.Path, entries.Name)
	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
