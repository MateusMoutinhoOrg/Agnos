package remove_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_adapter"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeAdapterAction.RemoveAdapter(deps, entries.Path, entries.Adapter)

	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
