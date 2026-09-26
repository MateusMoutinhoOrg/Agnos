package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	// --priority and --segments declare no default, so a value typed as 0
	// is told apart from none by whether anything was bound at all.
	set_error := setRouteAction.SetRoute(sandbox, api.RouteProps{
		Priority:        command.GetInt("priority"),
		HasPriority:     command.GetItem("priority") != nil,
		Before:          command.GetString("before"),
		After:           command.GetString("after"),
		Segments:        command.GetInt("segments"),
		HasSegments:     command.GetItem("segments") != nil,
		Phase:           command.GetString("phase"),
		Clear:           command.GetStrings("clear"),
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
