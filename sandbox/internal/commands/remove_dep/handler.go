package remove_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_dep"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeDepAction.RemoveDep(sandbox, api.RemoveDepProps{
		Path:         command.GetString("path"),
		Dep:          command.GetString("dep"),
		WithAdapters: command.GetBool("with-adapters"),
	})

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
