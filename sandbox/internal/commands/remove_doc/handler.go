package remove_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_doc"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeDocAction.RemoveDoc(sandbox, command.GetString("path"), command.GetString("name"))
	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
