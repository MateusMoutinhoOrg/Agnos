package remove_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_parameter"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	err := removeParameterAction.RemoveParameter(sandbox, command.GetString("path"), command.GetString("route"), command.GetString("name"))
	if err != nil {
		sandbox.Deps.Std.Error("%s\n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
