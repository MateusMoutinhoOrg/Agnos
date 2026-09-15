package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_arg"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeArgAction.RemoveArg(sandbox, command.GetString("path"), command.GetString("command"), command.GetString("name"))
	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
