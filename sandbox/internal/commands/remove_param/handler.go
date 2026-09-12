package remove_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_param"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeParamAction.RemoveParam(sandbox, entries.Path, entries.Route, entries.Name)

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
