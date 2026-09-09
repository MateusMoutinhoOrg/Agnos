package remove_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_route"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeRouteAction.RemoveRoute(deps, entries.Path, entries.Name)

	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
