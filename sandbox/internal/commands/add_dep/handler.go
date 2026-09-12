package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	install_error := addDepAction.AddDep(sandbox, api.AddDepProps{
		Path:            entries.Path,
		Dep:             entries.Dep,
		Adapter:         entries.Adapter,
		As:              entries.As,
		RemoteAvailable: entries.RemoteAvailable,
	})

	if install_error != nil {
		sandbox.Deps.Std.Error("%s\n", install_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
