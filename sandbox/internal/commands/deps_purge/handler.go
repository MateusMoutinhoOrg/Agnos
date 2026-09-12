package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_purge"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	purge_error := depsPurgeAction.DepsPurge(sandbox, entries.Path)

	if purge_error != nil {
		sandbox.Deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
