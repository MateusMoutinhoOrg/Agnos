package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addCommandAction.AddCommand(sandbox, entries.Path, entries.Name, entries.Help, entries.Category)

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
