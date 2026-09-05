package add_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_cli_example"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	add_error := addCliExampleAction.AddCliExample(deps, entries.Path, entries.Name)
	if add_error != nil {
		deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
