package add_database_migration

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDatabaseMigrationAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_database_migration"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addDatabaseMigrationAction.AddDatabaseMigration(sandbox, api.AddDatabaseMigrationProps{
		Path:     props.Path,
		Database: input.Database,
	})

	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
