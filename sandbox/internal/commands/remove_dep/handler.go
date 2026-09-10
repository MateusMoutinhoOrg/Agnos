package remove_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_dep"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeDepAction.RemoveDep(deps, api.RemoveDepProps{
		Path:         entries.Path,
		Dep:          entries.Dep,
		WithAdapters: entries.WithAdapters,
	})

	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
