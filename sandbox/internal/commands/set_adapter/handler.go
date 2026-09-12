package set_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_adapter"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	set_error := setAdapterAction.SetAdapter(sandbox, api.SetAdapterProps{
		Path:      entries.Path,
		Dep:       entries.Dep,
		Adapter:   entries.Adapter,
		Available: entries.Available,
	})

	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
