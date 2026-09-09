package remove_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removePageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_page"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removePageAction.RemovePage(deps, entries.Path, entries.Name)

	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
