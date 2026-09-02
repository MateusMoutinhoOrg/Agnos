package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func DepsPurge(deps *deps.Deps, path string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := DepsPurgeInternal(deps, io, path); err != nil {
		return err
	}
	if err := buildAction.BuildInternal(deps, io, path); err != nil {
		return err
	}
	return io.Persist()
}
