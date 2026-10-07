package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	purge_error := depsPurgeAction.DepsPurge(sandbox, props.Path)

	if purge_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", purge_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
