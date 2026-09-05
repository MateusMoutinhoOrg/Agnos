package remove_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_lib_example"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeLibExampleAction.RemoveLibExample(deps, entries.Path, entries.Name)
	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
