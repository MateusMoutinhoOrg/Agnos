package show_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ShowDatabase renders one database's whole declaration as the lines of a
// tree. Like the list actions it opens a SmartIO and never calls io.Persist,
// and it runs no follow-up build: reading a declaration changes nothing.
func ShowDatabase(sandbox *api.Sandbox, path string, database string) ([]string, error) {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	return ShowDatabaseInternal(sandbox, io, database)
}
