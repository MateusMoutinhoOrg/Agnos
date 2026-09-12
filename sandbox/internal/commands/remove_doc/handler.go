package remove_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_doc"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeDocAction.RemoveDoc(sandbox, entries.Path, entries.Name)
	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
