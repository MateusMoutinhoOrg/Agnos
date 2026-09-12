package server_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serverPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_purge"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	purge_error := serverPurgeAction.ServerPurge(sandbox, entries.Path)

	if purge_error != nil {
		sandbox.Deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
