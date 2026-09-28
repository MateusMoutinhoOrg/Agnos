package migrate_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	migrateCommandsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/migrate_commands"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	migrate_error := migrateCommandsAction.MigrateCommands(sandbox, props.Path)
	if migrate_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", migrate_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
