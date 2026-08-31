package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	startAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/start"
)

func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(path string) error {
		return buildAction.Build(deps, path)
	}
	sandbox.Actions.Start = func(props api.StartProps) error {
		return startAction.Start(deps, props)
	}
}
