package remove_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_doc"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeDocAction.RemoveDoc(deps, entries.Path, entries.Name)
	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
