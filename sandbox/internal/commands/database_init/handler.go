package database_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	databaseInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_init"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	init_error := databaseInitAction.DatabaseInit(sandbox, command.GetString("path"))

	if init_error != nil {
		sandbox.Deps.Std.Error("%s\n", init_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
