package set_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	set_error := setBodyAction.SetBody(sandbox, api.RouteBodyProps{
		Path:        command.GetString("path"),
		Route:       command.GetString("route"),
		Type:        command.GetString("type"),
		Required:    command.GetBool("required"),
		Optional:    command.GetBool("optional"),
		MaxBytes:    command.GetInt("max-bytes"),
		ContentType: command.GetString("content-type"),
		DropSchema:  command.GetBool("drop-schema"),
	})
	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
