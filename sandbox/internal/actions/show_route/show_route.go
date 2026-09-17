package show_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ShowRoute renders one route's whole declaration as the lines of a tree. Like
// the list actions it opens a SmartIO and never calls io.Persist, and it runs
// no follow-up build: reading a declaration changes nothing.
func ShowRoute(sandbox *api.Sandbox, path string, route string) ([]string, error) {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	return ShowRouteInternal(sandbox, io, route)
}
