package set_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// SetParameter rewrites one entry of the `parameters` of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a
// follow-up step. The build renders only: a renamed key may leave
// hand-written code reading an Input field that is gone.
func SetParameter(sandbox *api.Sandbox, props api.SetParameterProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := SetParameterInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
