package remove_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_available"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeAvailableAction.RemoveAvailable(sandbox, command.GetString("path"), command.GetString("available"))

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
