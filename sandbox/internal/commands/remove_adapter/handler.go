package remove_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_adapter"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeAdapterAction.RemoveAdapter(sandbox, entries.Path, entries.Adapter)

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
