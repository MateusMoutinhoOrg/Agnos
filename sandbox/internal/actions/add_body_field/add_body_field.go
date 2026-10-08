package add_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AddBodyField declares one property of the body json-schema of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a follow-up
// step so the Body struct, BodySchema and ReadBody pick it up.
func AddBodyField(sandbox *api.Sandbox, props api.AddBodyFieldProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddBodyFieldInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
