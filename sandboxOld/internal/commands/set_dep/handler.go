package set_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_dep"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	set_error := setDepAction.SetDep(sandbox, api.SetDepProps{
		Path:            command.GetString("path"),
		Dep:             command.GetString("dep"),
		Version:         command.GetString("version"),
		RemoteAvailable: command.GetString("remote-available"),
	})

	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
