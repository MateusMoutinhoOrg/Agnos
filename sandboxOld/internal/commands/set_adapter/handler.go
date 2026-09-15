package set_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_adapter"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	set_error := setAdapterAction.SetAdapter(sandbox, api.SetAdapterProps{
		Path:      command.GetString("path"),
		Dep:       command.GetString("dep"),
		Adapter:   command.GetString("adapter"),
		Available: command.GetString("available"),
	})

	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
