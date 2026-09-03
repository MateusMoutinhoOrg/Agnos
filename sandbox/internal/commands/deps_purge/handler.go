package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/deps_purge"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	purge_error := depsPurgeAction.DepsPurge(deps, entries.Path)

	if purge_error != nil {
		deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
