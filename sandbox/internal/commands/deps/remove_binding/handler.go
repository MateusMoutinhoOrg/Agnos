package remove_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_available"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeAvailableAction.RemoveAvailable(sandbox, props.Path, entries.Available)

	if remove_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
