package list_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ListRoutes renders the chain in the order the dispatch runs it. Like the
// other readers it opens a SmartIO and never calls io.Persist, and it runs no
// follow-up build.
func ListRoutes(sandbox *api.Sandbox, path string) ([]string, error) {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	return ListRoutesInternal(sandbox, io)
}
