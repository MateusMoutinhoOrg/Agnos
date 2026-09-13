package add_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_available"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addAvailableAction.AddAvailable(sandbox, command.GetString("path"), command.GetString("available"))

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
