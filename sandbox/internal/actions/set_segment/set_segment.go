package set_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetSegment rewrites one declared segment of
// sandbox/internal/routes/<route>/route.yaml's path, then runs build as a
// follow-up step so the route's new.go and its match order pick the change up.
func SetSegment(sandbox *api.Sandbox, props api.RouteFieldEditProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := SetSegmentInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
