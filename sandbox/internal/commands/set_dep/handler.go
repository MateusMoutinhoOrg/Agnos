package set_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	setDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_dep"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	set_error := setDepAction.SetDep(deps, api.SetDepProps{
		Path:            entries.Path,
		Dep:             entries.Dep,
		Version:         entries.Version,
		RemoteAvailable: entries.RemoteAvailable,
	})

	if set_error != nil {
		deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
