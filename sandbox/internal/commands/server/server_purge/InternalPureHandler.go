package server_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serverPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	purge_error := serverPurgeAction.ServerPurge(sandbox, props.Path)

	if purge_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", purge_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
