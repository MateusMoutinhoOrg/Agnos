package add_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table_field"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addTableFieldAction.AddTableField(sandbox, api.DatabaseFieldProps{
		Path:     command.GetString("path"),
		Database: command.GetString("database"),
		Table:    command.GetString("table"),
		Parent:   command.GetString("parent"),
		Name:     command.GetString("name"),
		Type:     command.GetString("type"),
		Required: command.GetBool("required"),
		Target:   command.GetString("target"),
	})

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
