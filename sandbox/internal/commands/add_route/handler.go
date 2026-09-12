package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addRouteAction.AddRoute(sandbox, entries.Path, entries.Name, entries.Method, entries.Trigger, entries.Help, entries.Category)

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
