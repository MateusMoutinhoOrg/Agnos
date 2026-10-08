package remove_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// RemoveParameter deletes one entry of the `parameters` of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a
// follow-up step. The build renders only: dropping a parameter may leave
// hand-written code reading an Input field that is gone.
func RemoveParameter(sandbox *api.Sandbox, props api.RemoveParameterProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RemoveParameterInternal(sandbox, io, props.Route, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
