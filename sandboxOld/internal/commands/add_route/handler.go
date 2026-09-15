package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addRouteAction.AddRoute(sandbox, command.GetString("path"), command.GetString("name"), command.GetString("method"), command.GetString("trigger"), command.GetString("help"), command.GetString("category"))

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
