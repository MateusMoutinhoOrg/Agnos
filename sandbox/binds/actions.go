package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	startAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/start"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(path string) error {
		io := smartio.New(deps, path, config.ProjectName)
		err := buildAction.Build(deps, io, path)
		if err != nil {
			return err
		}
		return io.Persist()
	}
	sandbox.Actions.Start = func(props api.StartProps) error {
		io := smartio.New(deps, props.Path, props.ProjectName)
		err := startAction.Start(deps, io, props)
		if err != nil {
			return err
		}
		return io.Persist()
	}
}
