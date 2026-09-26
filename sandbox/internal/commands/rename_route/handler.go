package rename_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	renameRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rename_route"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	rename_error := renameRouteAction.RenameRoute(sandbox, api.RenameRouteProps{
		Path:  command.GetString("path"),
		Route: command.GetString("route"),
		Name:  command.GetString("name"),
	})
	if rename_error != nil {
		sandbox.Deps.Std.Error("%s\n", rename_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
