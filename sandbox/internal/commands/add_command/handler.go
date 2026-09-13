package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addCommandAction.AddCommand(sandbox, command.GetString("path"), command.GetString("name"), command.GetString("help"), command.GetString("category"))

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
