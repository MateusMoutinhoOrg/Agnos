package add_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddTableField appends one field to a table of
// sandbox/internal/databases/<database>/specs.yaml, then runs build as a
// follow-up step so the methods that field generates are written.
func AddTableField(sandbox *api.Sandbox, props api.DatabaseFieldProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddTableFieldInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
