package remove_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_body_field"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeBodyFieldAction.RemoveBodyField(sandbox, entries.Path, entries.Route, entries.Name)

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
