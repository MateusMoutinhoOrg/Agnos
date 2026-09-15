package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetRoute rewrites the route-level keys of one route's route.yaml, then runs
// build as a follow-up step.
func SetRoute(sandbox *api.Sandbox, props api.RouteProps) error {
	io := smartio.New(sandbox, props.Path, config.ProjectName)
	if err := SetRouteInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
