package cli_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	cliPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_purge"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	purge_error := cliPurgeAction.CliPurge(sandbox, entries.Path)

	if purge_error != nil {
		sandbox.Deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
