package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddRoute scaffolds a new route package under
// sandbox/internal/routes/<name>/ — a declared route.yaml and a stub
// handler.go — then runs build as a follow-up step so entries.go and the
// dispatch arm are generated for it.
func AddRoute(deps *deps.Deps, path string, name string, method string, trigger string, help string, category string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := AddRouteInternal(deps, io, name, method, trigger, help, category); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
