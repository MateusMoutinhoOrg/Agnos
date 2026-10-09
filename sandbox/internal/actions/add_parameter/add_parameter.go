package add_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AddParameter inserts one entry into the `parameters` of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a
// follow-up step so the route's generated.new.go and generated.input.go pick it up.
func AddParameter(sandbox *api.Sandbox, props api.AddParameterProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddParameterInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
