package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_command"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeCommandAction.RemoveCommand(sandbox, entries.Path, entries.Name)
	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
