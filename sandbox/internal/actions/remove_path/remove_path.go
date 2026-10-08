package remove_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// RemovePath deletes one entry of the `paths` of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a
// follow-up step. The build renders only: dropping a path may leave
// hand-written code reading an Input field that is gone.
func RemovePath(sandbox *api.Sandbox, props api.RemovePathProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RemovePathInternal(sandbox, io, props.Route, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
