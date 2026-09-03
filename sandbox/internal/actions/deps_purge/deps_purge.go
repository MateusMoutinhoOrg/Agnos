package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func DepsPurge(deps *deps.Deps, path string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := DepsPurgeInternal(deps, io, path); err != nil {
		return err
	}
	if err := buildAction.BuildInternal(deps, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.RunRuntime(deps, path, api.RuntimeNone)
}
