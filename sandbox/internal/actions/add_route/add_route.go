package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AddRoute scaffolds a new route package under
// sandbox/internal/routes/<name>/ — a declared route.yaml and a stub
// handler.go — then runs build as a follow-up step so its new.go —
// the api.Route that lands in Server.Routes — and its input.go are
// generated for it.
func AddRoute(sandbox *api.Sandbox, props api.AddRouteProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddRouteInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
