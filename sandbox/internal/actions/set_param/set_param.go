package set_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetParam rewrites one declared query parameter of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a follow-up
// step so the route's new.go picks the change up.
func SetParam(sandbox *api.Sandbox, props api.RouteFieldEditProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := SetParamInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
