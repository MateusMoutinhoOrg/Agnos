package remove_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_flag"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeFlagAction.RemoveFlag(sandbox, command.GetString("path"), command.GetString("command"), command.GetString("name"))
	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
