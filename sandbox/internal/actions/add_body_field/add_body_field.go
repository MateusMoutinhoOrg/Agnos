package add_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddBodyField declares one property of the body json-schema of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a follow-up
// step so the Body struct, EntriesSchema and ReadBody pick it up.
func AddBodyField(sandbox *api.Sandbox, props api.RouteBodyFieldProps) error {
	io := smartio.New(sandbox, props.Path, config.ProjectName)
	if err := AddBodyFieldInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
