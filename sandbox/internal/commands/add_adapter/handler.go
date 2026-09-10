package add_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	install_error := addAdapterAction.AddAdapter(deps, api.AddAdapterProps{
		Path:      entries.Path,
		Adapter:   entries.Adapter,
		Available: entries.Available,
	})

	if install_error != nil {
		deps.Std.Error("%s\n", install_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
