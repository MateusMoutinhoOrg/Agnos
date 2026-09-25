package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addRouteAction.AddRoute(sandbox, api.AddRouteProps{
		Path:         command.GetString("path"),
		Name:         command.GetString("name"),
		Methods:      command.GetStrings("method"),
		Trigger:      command.GetString("trigger"),
		TriggerType:  command.GetString("trigger-type"),
		Priority:     command.GetInt("priority"),
		ResponseType: command.GetString("response-type"),
		Help:         command.GetString("help"),
		Category:     command.GetString("category"),
	})

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
