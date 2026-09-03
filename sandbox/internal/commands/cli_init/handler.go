package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	cliInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_init"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	init_error := cliInitAction.CliInit(deps, entries.Path)

	if init_error != nil {
		deps.Std.Error("%s\n", init_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
