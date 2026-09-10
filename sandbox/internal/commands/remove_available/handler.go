package remove_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_available"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeAvailableAction.RemoveAvailable(deps, entries.Path, entries.Available)

	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
