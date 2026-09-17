package set_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_param"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	set_error := setParamAction.SetParam(sandbox, api.RouteFieldEditProps{
		Path:        command.GetString("path"),
		Route:       command.GetString("route"),
		Name:        command.GetString("name"),
		Rename:      command.GetString("rename"),
		Description: command.GetString("description"),
		Examples:    command.GetStrings("example"),
		Type:        command.GetString("type"),
		Default:     command.GetString("default"),
		Required:    command.GetBool("required"),
		Array:       command.GetBool("array"),
		Min:         command.GetString("min"),
		Max:         command.GetString("max"),
		Clear:       command.GetStrings("clear"),
	})
	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
