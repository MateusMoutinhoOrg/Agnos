package add_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AddTable appends one table to
// sandbox/internal/databases/<database>/database.yaml, then runs build as a
// follow-up step so the package's three generated files pick it up.
func AddTable(sandbox *api.Sandbox, props api.AddTableProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddTableInternal(sandbox, io, props.Database, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
