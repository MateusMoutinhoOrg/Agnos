package add_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_available"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addAvailableAction.AddAvailable(sandbox, entries.Path, entries.Available)

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
