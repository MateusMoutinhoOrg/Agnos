package import_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ImportBody declares the body json-schema of
// sandbox/internal/routeslist/<route>/route.yaml from an example payload, then
// runs build as a follow-up step so the Body struct, BodySchema and ReadBody
// pick the properties up.
func ImportBody(sandbox *api.Sandbox, props api.RouteBodyImportProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := ImportBodyInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
