package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_purge"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	purge_error := depsPurgeAction.DepsPurge(sandbox, command.GetString("path"))

	if purge_error != nil {
		sandbox.Deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
