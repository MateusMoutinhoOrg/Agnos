package remove_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_field"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeFieldAction.RemoveField(deps, entries.Path, entries.Route, entries.In, entries.Name)

	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
