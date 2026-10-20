package add_database_migration

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AddDatabaseMigration opens the next version of one database of
// sandbox/internal/databases/<database>/database.yaml and writes its
// hand-written step, then runs build as a follow-up step so the database's
// generated.migrations.go runs it.
func AddDatabaseMigration(sandbox *api.Sandbox, props api.AddDatabaseMigrationProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddDatabaseMigrationInternal(sandbox, io, props.Database); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
