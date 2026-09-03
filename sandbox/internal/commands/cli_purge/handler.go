package cli_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	cliPurgeAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/cli_purge"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	purge_error := cliPurgeAction.CliPurge(deps, entries.Path)

	if purge_error != nil {
		deps.Std.Error("%s\n", purge_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
