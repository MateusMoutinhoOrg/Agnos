package front_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	init_error := frontInitAction.FrontInit(sandbox, command.GetString("path"))

	if init_error != nil {
		sandbox.Deps.Std.Error("%s\n", init_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
