package remove_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// RemoveTable drops one table from
// sandbox/internal/databases/<database>/database.yaml, then runs build as a
// follow-up step. The build renders only: dropping a table takes every method
// spelled after it with it, and hand-written code may still be calling one.
func RemoveTable(sandbox *api.Sandbox, props api.RemoveTableProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RemoveTableInternal(sandbox, io, props.Database, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
