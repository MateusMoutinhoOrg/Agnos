package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_arg"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeArgAction.RemoveArg(sandbox, entries.Path, entries.Command, entries.Name)
	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
