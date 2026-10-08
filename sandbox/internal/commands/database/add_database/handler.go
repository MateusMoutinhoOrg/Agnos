package add_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_database"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addDatabaseAction.AddDatabase(sandbox, api.AddDatabaseProps{Path: props.Path, Name: input.Name, KeyPrefix: input.KeyPrefix})

	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
