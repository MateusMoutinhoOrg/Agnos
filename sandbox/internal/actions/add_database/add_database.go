package add_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddDatabase writes the declaration of a new database under
// sandbox/internal/databases/<name>/specs.yaml, then runs build as a follow-up
// step so its api.go, new.go and methods.go are generated for it.
func AddDatabase(sandbox *api.Sandbox, path string, name string, prefix string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := AddDatabaseInternal(sandbox, io, name, prefix); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
