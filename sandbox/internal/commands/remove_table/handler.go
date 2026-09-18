package remove_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeTableAction.RemoveTable(
		sandbox,
		command.GetString("path"),
		command.GetString("database"),
		command.GetString("name"),
	)

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
