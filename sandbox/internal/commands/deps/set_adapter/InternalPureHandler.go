package set_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_adapter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	set_error := setAdapterAction.SetAdapter(sandbox, api.SetAdapterProps{
		Path:      props.Path,
		Dep:       entries.Dep,
		Adapter:   entries.Adapter,
		Available: entries.Available,
	})

	if set_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
