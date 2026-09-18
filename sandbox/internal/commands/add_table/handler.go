package add_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addTableAction.AddTable(
		sandbox,
		command.GetString("path"),
		command.GetString("database"),
		command.GetString("name"),
	)

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
