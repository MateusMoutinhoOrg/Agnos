package cli_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	cliPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	purge_error := cliPurgeAction.CliPurge(sandbox, props.Path)

	if purge_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", purge_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
