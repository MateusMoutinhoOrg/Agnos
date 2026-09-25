package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	// --priority defaults to -1, the one value the declaration refuses, so
	// setting the rung to 0 is told apart from not naming one at all.
	priority := command.GetInt("priority")

	set_error := setRouteAction.SetRoute(sandbox, api.RouteProps{
		Priority:        priority,
		HasPriority:     priority >= 0,
		Path:            command.GetString("path"),
		Route:           command.GetString("route"),
		Methods:         command.GetStrings("method"),
		ResponseType:    command.GetString("response-type"),
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
