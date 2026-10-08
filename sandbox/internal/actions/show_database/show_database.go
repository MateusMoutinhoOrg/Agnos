package show_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// ShowDatabase renders one database's whole declaration as the lines of a
// tree. Like the list actions it opens a StagedFS and never calls io.Persist,
// and it runs no follow-up build: reading a declaration changes nothing.
func ShowDatabase(sandbox *api.Sandbox, props api.ShowDatabaseProps) ([]string, error) {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return ShowDatabaseInternal(sandbox, io, props.Name)
}
