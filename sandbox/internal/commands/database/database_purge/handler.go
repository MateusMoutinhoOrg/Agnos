package database_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	databasePurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	purge_error := databasePurgeAction.DatabasePurge(sandbox, api.DatabasePurgeProps{Path: props.Path})

	if purge_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", purge_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
