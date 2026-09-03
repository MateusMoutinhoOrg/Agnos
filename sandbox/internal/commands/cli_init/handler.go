package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	cliInitAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/cli_init"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	init_error := cliInitAction.CliInit(deps, entries.Path)

	if !entries.Quiet && init_error != nil {
		deps.Std.Error(init_error.Error())
	}
	if init_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
