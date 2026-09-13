package add_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addHeaderAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_header"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addHeaderAction.AddHeader(sandbox, api.RouteFieldProps{
		Path:        command.GetString("path"),
		Route:       command.GetString("route"),
		Name:        command.GetString("name"),
		Description: command.GetString("description"),
		Examples:    command.GetStrings("example"),
		Type:        command.GetString("type"),
		Default:     command.GetString("default"),
		Required:    command.GetBool("required"),
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
