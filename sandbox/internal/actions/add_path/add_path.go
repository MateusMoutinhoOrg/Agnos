package add_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddPath inserts one entry into the `paths` of
// sandbox/internal/routeslist/<route>/route.yaml, then runs build as a
// follow-up step so the route's new.go and entries.go pick it up.
func AddPath(sandbox *api.Sandbox, props api.RoutePathProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddPathInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
