package remove_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemovePath deletes one entry of the `paths` of
// sandbox/internal/routeslist/<route>/route.yaml, then runs build as a
// follow-up step. The build renders only: dropping a path may leave
// hand-written code reading an Entries field that is gone.
func RemovePath(sandbox *api.Sandbox, path string, route string, id string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := RemovePathInternal(sandbox, io, route, id); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
