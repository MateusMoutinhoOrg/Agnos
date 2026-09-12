package remove_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_cli_example"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeCliExampleAction.RemoveCliExample(sandbox, entries.Path, entries.Name)
	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
