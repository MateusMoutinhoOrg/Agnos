package add_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addPathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_path"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	err := addPathAction.AddPath(sandbox, api.RoutePathProps{
		Path:              command.GetString("path"),
		Route:             command.GetString("route"),
		Id:                command.GetString("id"),
		Start:             command.GetString("start"),
		End:               command.GetString("end"),
		TriggerType:       command.GetString("trigger-type"),
		TriggerNegate:     command.GetBool("trigger-negate"),
		TriggerIgnoreCase: command.GetBool("trigger-ignore-case"),
		Type:              command.GetString("type"),
		Trigger:           command.GetString("trigger"),
		Description:       command.GetString("description"),
		Position:          command.GetInt("position"),
	})
	if err != nil {
		sandbox.Deps.Std.Error("%s\n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
