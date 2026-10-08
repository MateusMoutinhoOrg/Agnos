package add_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	install_error := addAdapterAction.AddAdapter(sandbox, api.AddAdapterProps{
		Path:      props.Path,
		Adapter:   entries.Adapter,
		Available: entries.Available,
	})

	if install_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", install_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
