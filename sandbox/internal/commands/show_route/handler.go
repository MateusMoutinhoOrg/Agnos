package show_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	showRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_route"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	lines, show_error := showRouteAction.ShowRoute(sandbox, command.GetString("path"), command.GetString("route"))

	if show_error != nil {
		sandbox.Deps.Std.Error("%s\n", show_error.Error())
		return api.ExitFailure
	}

	for _, line := range lines {
		sandbox.Deps.Std.Printf("%s\n", line)
	}
	return api.ExitOk
}
