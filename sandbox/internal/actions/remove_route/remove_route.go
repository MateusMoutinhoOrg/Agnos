package remove_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveRoute deletes one route package whole, then runs build as a follow-up
// step so the dispatch stops naming it. The build renders only: dropping a
// route may leave hand-written code referring to what is gone.
func RemoveRoute(deps *deps.Deps, path string, name string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := RemoveRouteInternal(deps, io, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
