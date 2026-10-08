package add_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AddDatabase writes the declaration of a new database under
// sandbox/internal/databases/<name>/database.yaml, then runs build as a follow-up
// step so its api.go, new.go and methods.go are generated for it.
func AddDatabase(sandbox *api.Sandbox, props api.AddDatabaseProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddDatabaseInternal(sandbox, io, props.Name, props.KeyPrefix); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
