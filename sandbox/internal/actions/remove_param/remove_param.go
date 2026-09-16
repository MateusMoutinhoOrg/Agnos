package remove_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveParam deletes one declared query parameter from
// sandbox/internal/routes/<route>/route.yaml, then runs build as a follow-up
// step. The build renders only: dropping a parameter may leave hand-written
// code referring to what is gone.
func RemoveParam(sandbox *api.Sandbox, path string, route string, name string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := RemoveParamInternal(sandbox, io, route, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
