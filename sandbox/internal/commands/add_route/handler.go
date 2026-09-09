package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	add_error := addRouteAction.AddRoute(deps, entries.Path, entries.Name, entries.Method, entries.Trigger, entries.Help, entries.Category)

	if add_error != nil {
		deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
