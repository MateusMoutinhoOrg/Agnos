package set_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetBodyField rewrites one property of the body json-schema of
// sandbox/internal/routeslist/<route>/route.yaml, then runs build as a follow-up
// step so the Body struct, BodySchema and ReadBody pick the change up.
func SetBodyField(sandbox *api.Sandbox, props api.RouteBodyFieldEditProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := SetBodyFieldInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
