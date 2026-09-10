package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_arg"
	addBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_body_field"
	addCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_cli_example"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
	addDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_doc"
	addFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_flag"
	addHeaderAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_header"
	addLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_lib_example"
	addPageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_page"
	addParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_param"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
	addSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_segment"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	cliInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_init"
	cliPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_purge"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	listDepsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_deps"
	removeDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_dep"
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_init"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_purge"
	execTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/exec_tests"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
	frontPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_purge"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_arg"
	removeBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_body_field"
	removeCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_cli_example"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_command"
	removeDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_doc"
	removeFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_flag"
	removeHeaderAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_header"
	removeLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_lib_example"
	removePageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_page"
	removeParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_param"
	removeRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_route"
	removeSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_segment"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	serverPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_purge"
	setBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body"
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
	sandbox.Actions.AddDep = func(props api.AddDepProps) error {
		return addDepAction.AddDep(deps, props)
	}
	sandbox.Actions.RemoveDep = func(path string, dep string) error {
		return removeDepAction.RemoveDep(deps, path, dep)
	}
	sandbox.Actions.ListDeps = func(path string) ([]string, error) {
		return listDepsAction.ListDeps(deps, path)
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
	sandbox.Actions.AddSegment = func(props api.RouteFieldProps) error {
		return addSegmentAction.AddSegment(deps, props)
	}
	sandbox.Actions.RemoveSegment = func(path string, route string, name string) error {
		return removeSegmentAction.RemoveSegment(deps, path, route, name)
	}
	sandbox.Actions.AddHeader = func(props api.RouteFieldProps) error {
		return addHeaderAction.AddHeader(deps, props)
	}
	sandbox.Actions.RemoveHeader = func(path string, route string, name string) error {
		return removeHeaderAction.RemoveHeader(deps, path, route, name)
	}
	sandbox.Actions.AddParam = func(props api.RouteFieldProps) error {
		return addParamAction.AddParam(deps, props)
	}
	sandbox.Actions.RemoveParam = func(path string, route string, name string) error {
		return removeParamAction.RemoveParam(deps, path, route, name)
	}
	sandbox.Actions.SetBody = func(props api.RouteBodyProps) error {
		return setBodyAction.SetBody(deps, props)
	}
	sandbox.Actions.AddBodyField = func(props api.RouteBodyFieldProps) error {
		return addBodyFieldAction.AddBodyField(deps, props)
	}
	sandbox.Actions.RemoveBodyField = func(path string, route string, name string) error {
		return removeBodyFieldAction.RemoveBodyField(deps, path, route, name)
	}
	sandbox.Actions.FrontInit = func(path string) error {
		return frontInitAction.FrontInit(deps, path)
	}
	sandbox.Actions.FrontPurge = func(path string) error {
		return frontPurgeAction.FrontPurge(deps, path)
	}
	sandbox.Actions.AddPage = func(props api.PageProps) error {
		return addPageAction.AddPage(deps, props)
	}
	sandbox.Actions.RemovePage = func(path string, name string) error {
		return removePageAction.RemovePage(deps, path, name)
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
