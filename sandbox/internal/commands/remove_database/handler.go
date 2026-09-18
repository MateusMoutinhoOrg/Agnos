package remove_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_database"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeDatabaseAction.RemoveDatabase(sandbox, command.GetString("path"), command.GetString("name"))

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
