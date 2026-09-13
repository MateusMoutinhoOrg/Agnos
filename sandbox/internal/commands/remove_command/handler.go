package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_command"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeCommandAction.RemoveCommand(sandbox, command.GetString("path"), command.GetString("name"))
	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
