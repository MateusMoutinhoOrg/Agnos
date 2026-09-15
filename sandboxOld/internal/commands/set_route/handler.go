package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	set_error := setRouteAction.SetRoute(sandbox, api.RouteProps{
		Path:            command.GetString("path"),
		Route:           command.GetString("route"),
		Method:          command.GetString("method"),
		Help:            command.GetString("help"),
		Category:        command.GetString("category"),
		LongDescription: command.GetString("long-description"),
		Hidden:          command.GetBool("hidden"),
		Visible:         command.GetBool("visible"),
		Examples:        command.GetStrings("example"),
	})
	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
