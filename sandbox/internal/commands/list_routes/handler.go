package list_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listRoutesAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_routes"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	lines, list_error := listRoutesAction.ListRoutes(sandbox, command.GetString("path"))
	if list_error != nil {
		sandbox.Deps.Std.Error("%s\n", list_error.Error())
		return api.ExitFailure
	}

	for _, line := range lines {
		sandbox.Deps.Std.Printf("%s\n", line)
	}
	return api.ExitOk
}
