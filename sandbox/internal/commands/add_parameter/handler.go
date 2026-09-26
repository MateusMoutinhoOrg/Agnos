package add_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_parameter"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	err := addParameterAction.AddParameter(sandbox, api.RouteParameterProps{
		Path:              command.GetString("path"),
		Route:             command.GetString("route"),
		Name:              command.GetString("name"),
		Type:              command.GetString("type"),
		Fonts:             command.GetStrings("font"),
		Required:          command.GetBool("required"),
		Default:           command.GetString("default"),
		TriggerType:       command.GetString("trigger-type"),
		TriggerNegate:     command.GetBool("trigger-negate"),
		TriggerIgnoreCase: command.GetBool("trigger-ignore-case"),
		Trigger:           command.GetString("trigger"),
		Description:       command.GetString("description"),
		Examples:          command.GetStrings("example"),
		Position:          command.GetInt("position"),
	})
	if err != nil {
		sandbox.Deps.Std.Error("%s\n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
