package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	install_error := addDepAction.AddDep(sandbox, api.AddDepProps{
		Path:            command.GetString("path"),
		Dep:             command.GetString("dep"),
		Adapter:         command.GetString("adapter"),
		As:              command.GetString("as"),
		RemoteAvailable: command.GetString("remote-available"),
	})

	if install_error != nil {
		sandbox.Deps.Std.Error("%s\n", install_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
