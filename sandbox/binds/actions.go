package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	depInstallAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_install"
	depListAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_list"
	depRemoveAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_remove"
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/deps_init"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/deps_purge"
	startAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/start"
)

func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(path string) error {
		return buildAction.Build(deps, path)
	}
	sandbox.Actions.Start = func(props api.StartProps) error {
		return startAction.Start(deps, props)
	}
	sandbox.Actions.DepsInit = func(path string) error {
		return depsInitAction.DepsInit(deps, path)
	}
	sandbox.Actions.DepsPurge = func(path string) error {
		return depsPurgeAction.DepsPurge(deps, path)
	}
	sandbox.Actions.DepInstall = func(path string, dep string) error {
		return depInstallAction.DepInstall(deps, path, dep)
	}
	sandbox.Actions.DepRemove = func(path string, dep string) error {
		return depRemoveAction.DepRemove(deps, path, dep)
	}
	sandbox.Actions.DepList = func(path string) ([]string, error) {
		return depListAction.DepList(deps, path)
	}
}
