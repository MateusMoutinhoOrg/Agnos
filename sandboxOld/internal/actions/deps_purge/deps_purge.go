package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func DepsPurge(sandbox *api.Sandbox, path string) error {
	io := smartio.New(sandbox, path, config.ProjectName)
	if err := DepsPurgeInternal(sandbox, io, path); err != nil {
		return err
	}
	if err := buildAction.BuildInternal(sandbox, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.RunRuntime(sandbox, path, api.RuntimeNone)
}
