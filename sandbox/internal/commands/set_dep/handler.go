package set_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_dep"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	set_error := setDepAction.SetDep(sandbox, api.SetDepProps{
		Path:            entries.Path,
		Dep:             entries.Dep,
		Version:         entries.Version,
		RemoteAvailable: entries.RemoteAvailable,
	})

	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
