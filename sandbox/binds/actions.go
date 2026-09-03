package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/add_arg"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/add_command"
	addFlagAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/add_flag"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	cliInitAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/cli_init"
	cliPurgeAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/cli_purge"
	depInstallAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_install"
	depListAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_list"
	depRemoveAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_remove"
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/deps_init"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/deps_purge"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/remove_arg"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/remove_command"
	removeFlagAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/remove_flag"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/set_command"
	startAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/start"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/verify"
)

func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(path string) error {
		return buildAction.Build(deps, path)
	}
	sandbox.Actions.Verify = func(path string) error {
		return verifyAction.Verify(deps, path)
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
	sandbox.Actions.CliInit = func(path string) error {
		return cliInitAction.CliInit(deps, path)
	}
	sandbox.Actions.CliPurge = func(path string) error {
		return cliPurgeAction.CliPurge(deps, path)
	}
	sandbox.Actions.AddCommand = func(path string, name string, help string, category string) error {
		return addCommandAction.AddCommand(deps, path, name, help, category)
	}
	sandbox.Actions.RemoveCommand = func(path string, name string) error {
		return removeCommandAction.RemoveCommand(deps, path, name)
	}
	sandbox.Actions.SetCommand = func(props api.CommandProps) error {
		return setCommandAction.SetCommand(deps, props)
	}
	sandbox.Actions.AddFlag = func(props api.FieldProps) error {
		return addFlagAction.AddFlag(deps, props)
	}
	sandbox.Actions.RemoveFlag = func(path string, command string, name string) error {
		return removeFlagAction.RemoveFlag(deps, path, command, name)
	}
	sandbox.Actions.AddArg = func(props api.FieldProps) error {
		return addArgAction.AddArg(deps, props)
	}
	sandbox.Actions.RemoveArg = func(path string, command string, name string) error {
		return removeArgAction.RemoveArg(deps, path, command, name)
	}
}
