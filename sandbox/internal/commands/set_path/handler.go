package set_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setPathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_path"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	err := setPathAction.SetPath(sandbox, api.RoutePathEditProps{
		Path:        command.GetString("path"),
		Route:       command.GetString("route"),
		Id:          command.GetString("id"),
		Rename:      command.GetString("rename"),
		Start:       command.GetString("start"),
		End:         command.GetString("end"),
		TriggerType: command.GetString("trigger-type"),
		Trigger:     command.GetString("trigger"),
		Description: command.GetString("description"),
		Clear:       command.GetStrings("clear"),
	})
	if err != nil {
		sandbox.Deps.Std.Error("%s\n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
