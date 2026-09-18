package add_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddTable appends one table to
// sandbox/internal/databases/<database>/specs.yaml, then runs build as a
// follow-up step so the package's three generated files pick it up.
func AddTable(sandbox *api.Sandbox, path string, database string, table string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := AddTableInternal(sandbox, io, database, table); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
