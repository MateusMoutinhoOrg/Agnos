package show_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	showDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_database"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	lines, show_error := showDatabaseAction.ShowDatabase(sandbox, command.GetString("path"), command.GetString("database"))

	if show_error != nil {
		sandbox.Deps.Std.Error("%s\n", show_error.Error())
		return api.ExitFailure
	}

	for _, line := range lines {
		sandbox.Deps.Std.Printf("%s\n", line)
	}
	return api.ExitOk
}
