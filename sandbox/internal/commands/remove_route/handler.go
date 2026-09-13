package remove_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_route"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeRouteAction.RemoveRoute(sandbox, command.GetString("path"), command.GetString("name"))

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
