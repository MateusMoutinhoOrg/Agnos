package migrate_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// MigrateCommands rewrites every command of the project from the
// entries.yaml + handler.go it was declared with into the command.yaml +
// InternalPureHandler.go the chain dispatch reads, and persists the result. It
// runs no build: the tree it leaves is the one the next build renders from.
func MigrateCommands(sandbox *api.Sandbox, path string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := MigrateCommandsInternal(sandbox, io); err != nil {
		return err
	}
	return io.Persist()
}
