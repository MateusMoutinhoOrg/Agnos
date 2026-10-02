package backoffice_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	backofficePurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/backoffice_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	purge_error := backofficePurgeAction.BackofficePurge(sandbox, props.Path)

	if purge_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", purge_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
