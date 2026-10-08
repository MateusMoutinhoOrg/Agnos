package show_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// ShowRoute renders one route's whole declaration as the lines of a tree. Like
// the list actions it opens a StagedFS and never calls io.Persist, and it runs
// no follow-up build: reading a declaration changes nothing.
func ShowRoute(sandbox *api.Sandbox, props api.ShowRouteProps) ([]string, error) {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return ShowRouteInternal(sandbox, io, props.Name)
}
