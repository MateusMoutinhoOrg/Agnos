package cli_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	cliPurgeAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/cli_purge"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	purge_error := cliPurgeAction.CliPurge(deps, entries.Path)

	if !entries.Quiet && purge_error != nil {
		deps.Std.Error(purge_error.Error())
	}
	if purge_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
