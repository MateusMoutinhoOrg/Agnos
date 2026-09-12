package remove_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_dep"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeDepAction.RemoveDep(sandbox, api.RemoveDepProps{
		Path:         entries.Path,
		Dep:          entries.Dep,
		WithAdapters: entries.WithAdapters,
	})

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
