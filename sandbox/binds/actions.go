package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_arg"
	addAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_available"
	addBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_body_field"
	addCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_cli_example"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
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
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_init"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_purge"
	execTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/exec_tests"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
	frontPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_purge"
	listAdaptersAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_adapters"
	listDepsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_deps"
	removeAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_adapter"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_arg"
	removeAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_available"
	removeBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_body_field"
	removeCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_cli_example"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_command"
	removeDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_dep"
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
	setAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_adapter"
	setBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_command"
	setDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_dep"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
	startAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/start"
	updateTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/update_tests"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/verify"
)

func ActionsBind(sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(props api.BuildProps) error {
		return buildAction.Build(sandbox, props)
	}
	sandbox.Actions.Compile = func(props api.CompileProps) error {
		return compileAction.Compile(sandbox, props)
	}
	sandbox.Actions.Verify = func(path string) error {
		return verifyAction.Verify(sandbox, path)
	}
	sandbox.Actions.Start = func(props api.StartProps) error {
		return startAction.Start(sandbox, props)
	}
	sandbox.Actions.DepsInit = func(path string) error {
		return depsInitAction.DepsInit(sandbox, path)
	}
	sandbox.Actions.DepsPurge = func(path string) error {
		return depsPurgeAction.DepsPurge(sandbox, path)
	}
	sandbox.Actions.AddDep = func(props api.AddDepProps) error {
		return addDepAction.AddDep(sandbox, props)
	}
	sandbox.Actions.RemoveDep = func(props api.RemoveDepProps) error {
		return removeDepAction.RemoveDep(sandbox, props)
	}
	sandbox.Actions.ListDeps = func(path string) ([]api.DepInfo, error) {
		return listDepsAction.ListDeps(sandbox, path)
	}
	sandbox.Actions.SetDep = func(props api.SetDepProps) error {
		return setDepAction.SetDep(sandbox, props)
	}
	sandbox.Actions.AddAdapter = func(props api.AddAdapterProps) error {
		return addAdapterAction.AddAdapter(sandbox, props)
	}
	sandbox.Actions.RemoveAdapter = func(path string, adapter string) error {
		return removeAdapterAction.RemoveAdapter(sandbox, path, adapter)
	}
	sandbox.Actions.SetAdapter = func(props api.SetAdapterProps) error {
		return setAdapterAction.SetAdapter(sandbox, props)
	}
	sandbox.Actions.ListAdapters = func(path string) ([]api.AdapterInfo, error) {
		return listAdaptersAction.ListAdapters(sandbox, path)
	}
	sandbox.Actions.AddAvailable = func(path string, available string) error {
		return addAvailableAction.AddAvailable(sandbox, path, available)
	}
	sandbox.Actions.RemoveAvailable = func(path string, available string) error {
		return removeAvailableAction.RemoveAvailable(sandbox, path, available)
	}
	sandbox.Actions.CliInit = func(path string) error {
		return cliInitAction.CliInit(sandbox, path)
	}
	sandbox.Actions.CliPurge = func(path string) error {
		return cliPurgeAction.CliPurge(sandbox, path)
	}
	sandbox.Actions.AddCommand = func(path string, name string, help string, category string) error {
		return addCommandAction.AddCommand(sandbox, path, name, help, category)
	}
	sandbox.Actions.RemoveCommand = func(path string, name string) error {
		return removeCommandAction.RemoveCommand(sandbox, path, name)
	}
	sandbox.Actions.SetCommand = func(props api.CommandProps) error {
		return setCommandAction.SetCommand(sandbox, props)
	}
	sandbox.Actions.AddFlag = func(props api.FieldProps) error {
		return addFlagAction.AddFlag(sandbox, props)
	}
	sandbox.Actions.RemoveFlag = func(path string, command string, name string) error {
		return removeFlagAction.RemoveFlag(sandbox, path, command, name)
	}
	sandbox.Actions.AddArg = func(props api.FieldProps) error {
		return addArgAction.AddArg(sandbox, props)
	}
	sandbox.Actions.RemoveArg = func(path string, command string, name string) error {
		return removeArgAction.RemoveArg(sandbox, path, command, name)
	}
	sandbox.Actions.ServerInit = func(path string) error {
		return serverInitAction.ServerInit(sandbox, path)
	}
	sandbox.Actions.ServerPurge = func(path string) error {
		return serverPurgeAction.ServerPurge(sandbox, path)
	}
	sandbox.Actions.AddRoute = func(path string, name string, method string, trigger string, help string, category string) error {
		return addRouteAction.AddRoute(sandbox, path, name, method, trigger, help, category)
	}
	sandbox.Actions.RemoveRoute = func(path string, name string) error {
		return removeRouteAction.RemoveRoute(sandbox, path, name)
	}
	sandbox.Actions.SetRoute = func(props api.RouteProps) error {
		return setRouteAction.SetRoute(sandbox, props)
	}
	sandbox.Actions.AddSegment = func(props api.RouteFieldProps) error {
		return addSegmentAction.AddSegment(sandbox, props)
	}
	sandbox.Actions.RemoveSegment = func(path string, route string, name string) error {
		return removeSegmentAction.RemoveSegment(sandbox, path, route, name)
	}
	sandbox.Actions.AddHeader = func(props api.RouteFieldProps) error {
		return addHeaderAction.AddHeader(sandbox, props)
	}
	sandbox.Actions.RemoveHeader = func(path string, route string, name string) error {
		return removeHeaderAction.RemoveHeader(sandbox, path, route, name)
	}
	sandbox.Actions.AddParam = func(props api.RouteFieldProps) error {
		return addParamAction.AddParam(sandbox, props)
	}
	sandbox.Actions.RemoveParam = func(path string, route string, name string) error {
		return removeParamAction.RemoveParam(sandbox, path, route, name)
	}
	sandbox.Actions.SetBody = func(props api.RouteBodyProps) error {
		return setBodyAction.SetBody(sandbox, props)
	}
	sandbox.Actions.AddBodyField = func(props api.RouteBodyFieldProps) error {
		return addBodyFieldAction.AddBodyField(sandbox, props)
	}
	sandbox.Actions.RemoveBodyField = func(path string, route string, name string) error {
		return removeBodyFieldAction.RemoveBodyField(sandbox, path, route, name)
	}
	sandbox.Actions.FrontInit = func(path string) error {
		return frontInitAction.FrontInit(sandbox, path)
	}
	sandbox.Actions.FrontPurge = func(path string) error {
		return frontPurgeAction.FrontPurge(sandbox, path)
	}
	sandbox.Actions.AddPage = func(props api.PageProps) error {
		return addPageAction.AddPage(sandbox, props)
	}
	sandbox.Actions.RemovePage = func(path string, name string) error {
		return removePageAction.RemovePage(sandbox, path, name)
	}
	sandbox.Actions.AddDoc = func(props api.DocProps) error {
		return addDocAction.AddDoc(sandbox, props)
	}
	sandbox.Actions.RemoveDoc = func(path string, name string) error {
		return removeDocAction.RemoveDoc(sandbox, path, name)
	}
	sandbox.Actions.AddCliExample = func(path string, name string) error {
		return addCliExampleAction.AddCliExample(sandbox, path, name)
	}
	sandbox.Actions.RemoveCliExample = func(path string, name string) error {
		return removeCliExampleAction.RemoveCliExample(sandbox, path, name)
	}
	sandbox.Actions.AddLibExample = func(path string, name string) error {
		return addLibExampleAction.AddLibExample(sandbox, path, name)
	}
	sandbox.Actions.RemoveLibExample = func(path string, name string) error {
		return removeLibExampleAction.RemoveLibExample(sandbox, path, name)
	}
	sandbox.Actions.ExecTest = func(props api.ExecTestProps) error {
		return execTestsAction.ExecTest(sandbox, props)
	}
	sandbox.Actions.UpdateTest = func(path string, name string) error {
		return updateTestsAction.UpdateTest(sandbox, path, name)
	}
}
