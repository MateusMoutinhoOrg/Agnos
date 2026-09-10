package add_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_available"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	add_error := addAvailableAction.AddAvailable(deps, entries.Path, entries.Available)

	if add_error != nil {
		deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
