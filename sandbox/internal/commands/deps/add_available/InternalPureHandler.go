package add_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_available"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addAvailableAction.AddAvailable(sandbox, props.Path, entries.Available)

	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
