package add_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_database"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addDatabaseAction.AddDatabase(
		sandbox,
		props.Path,
		entries.Name,
		entries.Prefix,
	)

	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
