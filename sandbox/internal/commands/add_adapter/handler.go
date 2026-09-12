package add_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	install_error := addAdapterAction.AddAdapter(sandbox, api.AddAdapterProps{
		Path:      entries.Path,
		Adapter:   entries.Adapter,
		Available: entries.Available,
	})

	if install_error != nil {
		sandbox.Deps.Std.Error("%s\n", install_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
