package deps_init

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/deps_init"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	init_error := depsInitAction.DepsInit(deps, entries.Path)

	if init_error != nil {
		deps.Std.Error("%s\n", init_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
