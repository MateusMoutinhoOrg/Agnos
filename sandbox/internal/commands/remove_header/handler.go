package remove_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeHeaderAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_header"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeHeaderAction.RemoveHeader(sandbox, entries.Path, entries.Route, entries.Name)

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
