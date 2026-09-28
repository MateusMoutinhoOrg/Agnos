package database_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	databaseInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	init_error := databaseInitAction.DatabaseInit(sandbox, props.Path)

	if init_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", init_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
