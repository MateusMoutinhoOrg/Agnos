package add_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AddPath inserts one entry into the `paths` of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a
// follow-up step so the route's generated.new.go and generated.input.go pick it up.
func AddPath(sandbox *api.Sandbox, props api.AddPathProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddPathInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
