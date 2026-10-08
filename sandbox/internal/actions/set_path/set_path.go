package set_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// SetPath rewrites one entry of the `paths` of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a
// follow-up step. The build renders only: a renamed id may leave hand-written
// code reading an Input field that is gone.
func SetPath(sandbox *api.Sandbox, props api.SetPathProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := SetPathInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
