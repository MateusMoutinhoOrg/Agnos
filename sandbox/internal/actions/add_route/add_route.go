package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddRoute scaffolds a new route package under
// sandbox/internal/routes/<name>/ — a declared route.yaml and a stub
// handler.go — then runs build as a follow-up step so its new.go — the
// api.Route that lands in Server.Routes — is generated for it.
func AddRoute(sandbox *api.Sandbox, props api.AddRouteProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddRouteInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
