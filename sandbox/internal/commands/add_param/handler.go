package add_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_param"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addParamAction.AddParam(sandbox, api.RouteFieldProps{
		Path:        command.GetString("path"),
		Route:       command.GetString("route"),
		Name:        command.GetString("name"),
		Description: command.GetString("description"),
		Examples:    command.GetStrings("example"),
		Type:        command.GetString("type"),
		Default:     command.GetString("default"),
		Required:    command.GetBool("required"),
		Array:       command.GetBool("array"),
		Min:         command.GetString("min"),
		Max:         command.GetString("max"),
		Position:    command.GetInt("position"),
	})
	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
