package remove_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_adapter"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeAdapterAction.RemoveAdapter(sandbox, command.GetString("path"), command.GetString("adapter"))

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
