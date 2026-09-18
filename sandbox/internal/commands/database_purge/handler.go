package database_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	databasePurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_purge"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	purge_error := databasePurgeAction.DatabasePurge(sandbox, command.GetString("path"))

	if purge_error != nil {
		sandbox.Deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
