package remove_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeHeaderAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_header"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeHeaderAction.RemoveHeader(deps, entries.Path, entries.Route, entries.Name)

	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
