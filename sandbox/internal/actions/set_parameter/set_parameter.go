package set_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetParameter rewrites one entry of the `parameters` of
// sandbox/internal/routeslist/<route>/route.yaml, then runs build as a
// follow-up step. The build renders only: a renamed key may leave
// hand-written code reading an Entries field that is gone.
func SetParameter(sandbox *api.Sandbox, props api.RouteParameterEditProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := SetParameterInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
