package remove_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_available"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeAvailableAction.RemoveAvailable(sandbox, props.Path, entries.Available)

	if remove_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
