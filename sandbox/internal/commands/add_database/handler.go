package add_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_database"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addDatabaseAction.AddDatabase(
		sandbox,
		command.GetString("path"),
		command.GetString("name"),
		command.GetString("prefix"),
	)

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
