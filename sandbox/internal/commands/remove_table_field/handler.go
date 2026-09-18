package remove_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table_field"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeTableFieldAction.RemoveTableField(sandbox, api.DatabaseFieldProps{
		Path:     command.GetString("path"),
		Database: command.GetString("database"),
		Table:    command.GetString("table"),
		Parent:   command.GetString("parent"),
		Name:     command.GetString("name"),
	})

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
