package front_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	frontPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_purge"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	purge_error := frontPurgeAction.FrontPurge(sandbox, entries.Path)

	if purge_error != nil {
		sandbox.Deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
