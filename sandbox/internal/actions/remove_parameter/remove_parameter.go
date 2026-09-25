package remove_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveParameter deletes one entry of the `parameters` of
// sandbox/internal/routeslist/<route>/route.yaml, then runs build as a
// follow-up step. The build renders only: dropping a parameter may leave
// hand-written code reading an Entries field that is gone.
func RemoveParameter(sandbox *api.Sandbox, path string, route string, name string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := RemoveParameterInternal(sandbox, io, route, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
