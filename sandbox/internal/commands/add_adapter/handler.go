package add_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	install_error := addAdapterAction.AddAdapter(sandbox, api.AddAdapterProps{
		Path:      command.GetString("path"),
		Adapter:   command.GetString("adapter"),
		Available: command.GetString("available"),
	})

	if install_error != nil {
		sandbox.Deps.Std.Error("%s\n", install_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
