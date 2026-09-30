package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	purge_error := depsPurgeAction.DepsPurge(sandbox, props.Path)

	if purge_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", purge_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
