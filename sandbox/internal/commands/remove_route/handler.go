package remove_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_route"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeRouteAction.RemoveRoute(sandbox, entries.Path, entries.Name)

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
