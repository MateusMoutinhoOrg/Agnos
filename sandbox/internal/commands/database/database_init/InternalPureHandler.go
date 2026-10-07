package database_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	databaseInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	init_error := databaseInitAction.DatabaseInit(sandbox, props.Path)

	if init_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", init_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
