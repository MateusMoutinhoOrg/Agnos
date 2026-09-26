package explain_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ExplainRoute runs one request against the declared routes without a server
// and renders what each one does with it. It writes nothing and runs no build.
func ExplainRoute(sandbox *api.Sandbox, props api.ExplainRouteProps) ([]string, error) {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return ExplainRouteInternal(sandbox, io, props)
}
