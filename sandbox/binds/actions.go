package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	enableDepsAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/enable_deps"
	removeDepsAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/remove_deps"
	startAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/start"
)

func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(path string) error {
		return buildAction.Build(deps, path)
	}
	sandbox.Actions.Start = func(props api.StartProps) error {
		return startAction.Start(deps, props)
	}
	sandbox.Actions.EnableDeps = func(path string) error {
		return enableDepsAction.EnableDeps(deps, path)
	}
	sandbox.Actions.RemoveDeps = func(path string) error {
		return removeDepsAction.RemoveDeps(deps, path)
	}
}
