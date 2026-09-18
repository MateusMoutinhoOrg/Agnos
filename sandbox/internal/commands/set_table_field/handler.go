package set_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_table_field"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	set_error := setTableFieldAction.SetTableField(sandbox, api.DatabaseFieldEditProps{
		Path:     command.GetString("path"),
		Database: command.GetString("database"),
		Table:    command.GetString("table"),
		Parent:   command.GetString("parent"),
		Name:     command.GetString("name"),
		Rename:   command.GetString("rename"),
		Type:     command.GetString("type"),
		Required: command.GetBool("required"),
		Target:   command.GetString("target"),
		Clear:    command.GetStrings("clear"),
	})

	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
