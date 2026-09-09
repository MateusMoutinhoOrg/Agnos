package front_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	frontPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_purge"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	purge_error := frontPurgeAction.FrontPurge(deps, entries.Path)

	if purge_error != nil {
		deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
