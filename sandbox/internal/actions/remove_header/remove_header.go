package remove_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveHeader deletes one declared header from
// sandbox/internal/routes/<route>/route.yaml, then runs build as a follow-up
// step. The build renders only: dropping a header may leave hand-written code
// referring to what is gone.
func RemoveHeader(deps *deps.Deps, path string, route string, name string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := RemoveHeaderInternal(deps, io, route, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
