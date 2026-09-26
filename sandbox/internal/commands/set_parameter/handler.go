package set_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_parameter"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	err := setParameterAction.SetParameter(sandbox, api.RouteParameterEditProps{
		Path:              command.GetString("path"),
		Route:             command.GetString("route"),
		Name:              command.GetString("name"),
		Rename:            command.GetString("rename"),
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
		Clear:             command.GetStrings("clear"),
	})
	if err != nil {
		sandbox.Deps.Std.Error("%s\n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
