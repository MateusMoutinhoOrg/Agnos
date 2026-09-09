package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	set_error := setRouteAction.SetRoute(deps, api.RouteProps{
		Path:            entries.Path,
		Route:           entries.Route,
		Method:          entries.Method,
		Help:            entries.Help,
		Category:        entries.Category,
		LongDescription: entries.LongDescription,
		Hidden:          entries.Hidden,
		Visible:         entries.Visible,
		Examples:        entries.Example,
	})
	if set_error != nil {
		deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
