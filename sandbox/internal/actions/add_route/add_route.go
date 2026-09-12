package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddRoute scaffolds a new route package under
// sandbox/internal/routes/<name>/ — a declared route.yaml and a stub
// handler.go — then runs build as a follow-up step so entries.go and the
// dispatch arm are generated for it.
func AddRoute(sandbox *api.Sandbox, path string, name string, method string, trigger string, help string, category string) error {
	io := smartio.New(sandbox, path, config.ProjectName)
	if err := AddRouteInternal(sandbox, io, name, method, trigger, help, category); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
