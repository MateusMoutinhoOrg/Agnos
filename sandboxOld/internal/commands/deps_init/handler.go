package deps_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_init"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	init_error := depsInitAction.DepsInit(sandbox, command.GetString("path"))

	if init_error != nil {
		sandbox.Deps.Std.Error("%s\n", init_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
