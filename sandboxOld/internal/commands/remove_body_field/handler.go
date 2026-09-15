package remove_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_body_field"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeBodyFieldAction.RemoveBodyField(sandbox, command.GetString("path"), command.GetString("route"), command.GetString("name"))

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
