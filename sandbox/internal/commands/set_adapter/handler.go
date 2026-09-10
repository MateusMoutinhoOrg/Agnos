package set_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	setAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_adapter"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	set_error := setAdapterAction.SetAdapter(deps, api.SetAdapterProps{
		Path:      entries.Path,
		Dep:       entries.Dep,
		Adapter:   entries.Adapter,
		Available: entries.Available,
	})

	if set_error != nil {
		deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
