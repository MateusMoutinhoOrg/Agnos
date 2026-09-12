package server_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	init_error := serverInitAction.ServerInit(sandbox, entries.Path)

	if init_error != nil {
		sandbox.Deps.Std.Error("%s\n", init_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
