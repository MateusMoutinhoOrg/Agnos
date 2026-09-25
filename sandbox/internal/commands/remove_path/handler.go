package remove_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removePathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_path"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	err := removePathAction.RemovePath(sandbox, command.GetString("path"), command.GetString("route"), command.GetString("id"))
	if err != nil {
		sandbox.Deps.Std.Error("%s\n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
