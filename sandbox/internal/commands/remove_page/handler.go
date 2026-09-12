package remove_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removePageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_page"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removePageAction.RemovePage(sandbox, entries.Path, entries.Name)

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
