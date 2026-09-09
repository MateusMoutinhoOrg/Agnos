package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_arg"
	addCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_cli_example"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
	addDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_doc"
	addFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_field"
	addFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_flag"
	addLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_lib_example"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	cliInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_init"
	cliPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_purge"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
	depInstallAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/dep_install"
	depListAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/dep_list"
	depRemoveAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/dep_remove"
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_init"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_purge"
	execTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/exec_tests"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_arg"
	removeCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_cli_example"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_command"
	removeDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_doc"
	removeFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_field"
	removeFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_flag"
	removeLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_lib_example"
	removeRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_route"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	serverPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_purge"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_command"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
	startAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/start"
	updateTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/update_tests"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/verify"
)

func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(props api.BuildProps) error {
		return buildAction.Build(deps, props)
	}
	sandbox.Actions.Compile = func(props api.CompileProps) error {
		return compileAction.Compile(deps, props)
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
	sandbox.Actions.ServerInit = func(path string) error {
		return serverInitAction.ServerInit(deps, path)
	}
	sandbox.Actions.ServerPurge = func(path string) error {
		return serverPurgeAction.ServerPurge(deps, path)
	}
	sandbox.Actions.AddRoute = func(path string, name string, method string, trigger string, help string, category string) error {
		return addRouteAction.AddRoute(deps, path, name, method, trigger, help, category)
	}
	sandbox.Actions.RemoveRoute = func(path string, name string) error {
		return removeRouteAction.RemoveRoute(deps, path, name)
	}
	sandbox.Actions.SetRoute = func(props api.RouteProps) error {
		return setRouteAction.SetRoute(deps, props)
	}
	sandbox.Actions.AddField = func(props api.RouteFieldProps) error {
		return addFieldAction.AddField(deps, props)
	}
	sandbox.Actions.RemoveField = func(path string, route string, in string, name string) error {
		return removeFieldAction.RemoveField(deps, path, route, in, name)
	}
	sandbox.Actions.AddDoc = func(props api.DocProps) error {
		return addDocAction.AddDoc(deps, props)
	}
	sandbox.Actions.RemoveDoc = func(path string, name string) error {
		return removeDocAction.RemoveDoc(deps, path, name)
	}
	sandbox.Actions.AddCliExample = func(path string, name string) error {
		return addCliExampleAction.AddCliExample(deps, path, name)
	}
	sandbox.Actions.RemoveCliExample = func(path string, name string) error {
		return removeCliExampleAction.RemoveCliExample(deps, path, name)
	}
	sandbox.Actions.AddLibExample = func(path string, name string) error {
		return addLibExampleAction.AddLibExample(deps, path, name)
	}
	sandbox.Actions.RemoveLibExample = func(path string, name string) error {
		return removeLibExampleAction.RemoveLibExample(deps, path, name)
	}
	sandbox.Actions.ExecTest = func(props api.ExecTestProps) error {
		return execTestsAction.ExecTest(deps, props)
	}
	sandbox.Actions.UpdateTest = func(path string, name string) error {
		return updateTestsAction.UpdateTest(deps, path, name)
	}
}
