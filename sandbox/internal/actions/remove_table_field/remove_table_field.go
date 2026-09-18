package remove_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveTableField drops one field from a table of
// sandbox/internal/databases/<database>/specs.yaml, then runs build as a
// follow-up step. The build renders only: dropping a field takes the methods
// it generated with it, and hand-written code may still be calling one.
func RemoveTableField(sandbox *api.Sandbox, props api.DatabaseFieldProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RemoveTableFieldInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
