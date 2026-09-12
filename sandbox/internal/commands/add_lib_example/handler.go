package add_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_lib_example"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addLibExampleAction.AddLibExample(sandbox, entries.Path, entries.Name)
	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
