package front_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	frontPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_purge"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	purge_error := frontPurgeAction.FrontPurge(sandbox, command.GetString("path"))

	if purge_error != nil {
		sandbox.Deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
