package add_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_cli_example"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addCliExampleAction.AddCliExample(sandbox, entries.Path, entries.Name)
	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
