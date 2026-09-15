package remove_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_param"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeParamAction.RemoveParam(sandbox, command.GetString("path"), command.GetString("route"), command.GetString("name"))

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
