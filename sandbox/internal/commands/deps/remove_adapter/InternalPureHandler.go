package remove_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_adapter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeAdapterAction.RemoveAdapter(sandbox, props.Path, entries.Adapter)

	if remove_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
