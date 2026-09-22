package actions

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_arg"
	addAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_available"
	addBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_body_field"
	addCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_cli_example"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
	addDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_database"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	addDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_doc"
	addFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_flag"
	addHeaderAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_header"
	addLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_lib_example"
	addPageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_page"
	addParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_param"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
	addSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_segment"
	addTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table"
	addTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table_field"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	cliInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_init"
	cliPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_purge"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
	databaseInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_init"
	databasePurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_purge"
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_init"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_purge"
	disableExtensionAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/disable_extension"
	enableExtensionAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/enable_extension"
	execTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/exec_tests"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
	frontPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_purge"
	importBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/import_body"
	interviewAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/interview"
	listAdaptersAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_adapters"
	listDepsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_deps"
	listExtensionsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_extensions"
	removeAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_adapter"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_arg"
	removeAvailableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_available"
	removeBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_body_field"
	removeCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_cli_example"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_command"
	removeDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_database"
	removeDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_dep"
	removeDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_doc"
	removeFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_flag"
	removeHeaderAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_header"
	removeLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_lib_example"
	removePageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_page"
	removeParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_param"
	removeRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_route"
	removeSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_segment"
	removeTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table"
	removeTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table_field"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	serverPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_purge"
	setAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_adapter"
	setBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body"
	setBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body_field"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_command"
	setDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_dep"
	setHeaderAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_header"
	setParamAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_param"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
	setSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_segment"
	setTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_table_field"
	showDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_database"
	showRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_route"
	startAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/start"
	updateTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/update_tests"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/verify"
)

// NewActions builds the action surface of the sandbox: one field of api.Actions
// per sandbox/internal/actions/<name>/, closed over the sandbox it was built
// from. Hand-written — a new action is added here alongside its field in
// sandbox/api/actions.go.
func NewActions(sandbox *api.Sandbox) api.Actions {
	actions := api.Actions{}

	actions.Build = func(props api.BuildProps) error {
		return buildAction.Build(sandbox, props)
	}
	actions.Compile = func(props api.CompileProps) error {
		return compileAction.Compile(sandbox, props)
	}
	actions.Verify = func(path string) error {
		return verifyAction.Verify(sandbox, path)
	}
	actions.Start = func(props api.StartProps) error {
		return startAction.Start(sandbox, props)
	}
	actions.EnableExtension = func(path string, name string) error {
		return enableExtensionAction.EnableExtension(sandbox, path, name)
	}
	actions.DisableExtension = func(path string, name string) error {
		return disableExtensionAction.DisableExtension(sandbox, path, name)
	}
	actions.ListExtensions = func(path string) ([]api.ExtensionInfo, error) {
		return listExtensionsAction.ListExtensions(sandbox, path)
	}
	actions.DepsInit = func(path string) error {
		return depsInitAction.DepsInit(sandbox, path)
	}
	actions.DepsPurge = func(path string) error {
		return depsPurgeAction.DepsPurge(sandbox, path)
	}
	actions.AddDep = func(props api.AddDepProps) error {
		return addDepAction.AddDep(sandbox, props)
	}
	actions.RemoveDep = func(props api.RemoveDepProps) error {
		return removeDepAction.RemoveDep(sandbox, props)
	}
	actions.ListDeps = func(path string) ([]api.DepInfo, error) {
		return listDepsAction.ListDeps(sandbox, path)
	}
	actions.SetDep = func(props api.SetDepProps) error {
		return setDepAction.SetDep(sandbox, props)
	}
	actions.AddAdapter = func(props api.AddAdapterProps) error {
		return addAdapterAction.AddAdapter(sandbox, props)
	}
	actions.RemoveAdapter = func(path string, adapter string) error {
		return removeAdapterAction.RemoveAdapter(sandbox, path, adapter)
	}
	actions.SetAdapter = func(props api.SetAdapterProps) error {
		return setAdapterAction.SetAdapter(sandbox, props)
	}
	actions.ListAdapters = func(path string) ([]api.AdapterInfo, error) {
		return listAdaptersAction.ListAdapters(sandbox, path)
	}
	actions.AddAvailable = func(path string, available string) error {
		return addAvailableAction.AddAvailable(sandbox, path, available)
	}
	actions.RemoveAvailable = func(path string, available string) error {
		return removeAvailableAction.RemoveAvailable(sandbox, path, available)
	}
	actions.CliInit = func(path string) error {
		return cliInitAction.CliInit(sandbox, path)
	}
	actions.CliPurge = func(path string) error {
		return cliPurgeAction.CliPurge(sandbox, path)
	}
	actions.AddCommand = func(path string, name string, help string, category string) error {
		return addCommandAction.AddCommand(sandbox, path, name, help, category)
	}
	actions.RemoveCommand = func(path string, name string) error {
		return removeCommandAction.RemoveCommand(sandbox, path, name)
	}
	actions.SetCommand = func(props api.CommandProps) error {
		return setCommandAction.SetCommand(sandbox, props)
	}
	actions.AddFlag = func(props api.FieldProps) error {
		return addFlagAction.AddFlag(sandbox, props)
	}
	actions.RemoveFlag = func(path string, command string, name string) error {
		return removeFlagAction.RemoveFlag(sandbox, path, command, name)
	}
	actions.AddArg = func(props api.FieldProps) error {
		return addArgAction.AddArg(sandbox, props)
	}
	actions.RemoveArg = func(path string, command string, name string) error {
		return removeArgAction.RemoveArg(sandbox, path, command, name)
	}
	actions.ServerInit = func(path string) error {
		return serverInitAction.ServerInit(sandbox, path)
	}
	actions.ServerPurge = func(path string) error {
		return serverPurgeAction.ServerPurge(sandbox, path)
	}
	actions.AddRoute = func(props api.AddRouteProps) error {
		return addRouteAction.AddRoute(sandbox, props)
	}
	actions.RemoveRoute = func(path string, name string) error {
		return removeRouteAction.RemoveRoute(sandbox, path, name)
	}
	actions.SetRoute = func(props api.RouteProps) error {
		return setRouteAction.SetRoute(sandbox, props)
	}
	actions.AddSegment = func(props api.RouteFieldProps) error {
		return addSegmentAction.AddSegment(sandbox, props)
	}
	actions.RemoveSegment = func(path string, route string, name string) error {
		return removeSegmentAction.RemoveSegment(sandbox, path, route, name)
	}
	actions.AddHeader = func(props api.RouteFieldProps) error {
		return addHeaderAction.AddHeader(sandbox, props)
	}
	actions.RemoveHeader = func(path string, route string, name string) error {
		return removeHeaderAction.RemoveHeader(sandbox, path, route, name)
	}
	actions.AddParam = func(props api.RouteFieldProps) error {
		return addParamAction.AddParam(sandbox, props)
	}
	actions.RemoveParam = func(path string, route string, name string) error {
		return removeParamAction.RemoveParam(sandbox, path, route, name)
	}
	actions.SetBody = func(props api.RouteBodyProps) error {
		return setBodyAction.SetBody(sandbox, props)
	}
	actions.AddBodyField = func(props api.RouteBodyFieldProps) error {
		return addBodyFieldAction.AddBodyField(sandbox, props)
	}
	actions.RemoveBodyField = func(path string, route string, name string) error {
		return removeBodyFieldAction.RemoveBodyField(sandbox, path, route, name)
	}
	actions.SetSegment = func(props api.RouteFieldEditProps) error {
		return setSegmentAction.SetSegment(sandbox, props)
	}
	actions.SetHeader = func(props api.RouteFieldEditProps) error {
		return setHeaderAction.SetHeader(sandbox, props)
	}
	actions.SetParam = func(props api.RouteFieldEditProps) error {
		return setParamAction.SetParam(sandbox, props)
	}
	actions.SetBodyField = func(props api.RouteBodyFieldEditProps) error {
		return setBodyFieldAction.SetBodyField(sandbox, props)
	}
	actions.ImportBody = func(props api.RouteBodyImportProps) error {
		return importBodyAction.ImportBody(sandbox, props)
	}
	actions.ShowRoute = func(path string, route string) ([]string, error) {
		return showRouteAction.ShowRoute(sandbox, path, route)
	}
	actions.DatabaseInit = func(path string) error {
		return databaseInitAction.DatabaseInit(sandbox, path)
	}
	actions.DatabasePurge = func(path string) error {
		return databasePurgeAction.DatabasePurge(sandbox, path)
	}
	actions.AddDatabase = func(path string, name string, prefix string) error {
		return addDatabaseAction.AddDatabase(sandbox, path, name, prefix)
	}
	actions.RemoveDatabase = func(path string, name string) error {
		return removeDatabaseAction.RemoveDatabase(sandbox, path, name)
	}
	actions.AddTable = func(path string, database string, table string) error {
		return addTableAction.AddTable(sandbox, path, database, table)
	}
	actions.RemoveTable = func(path string, database string, table string) error {
		return removeTableAction.RemoveTable(sandbox, path, database, table)
	}
	actions.AddTableField = func(props api.DatabaseFieldProps) error {
		return addTableFieldAction.AddTableField(sandbox, props)
	}
	actions.SetTableField = func(props api.DatabaseFieldEditProps) error {
		return setTableFieldAction.SetTableField(sandbox, props)
	}
	actions.RemoveTableField = func(props api.DatabaseFieldProps) error {
		return removeTableFieldAction.RemoveTableField(sandbox, props)
	}
	actions.ShowDatabase = func(path string, database string) ([]string, error) {
		return showDatabaseAction.ShowDatabase(sandbox, path, database)
	}
	actions.FrontInit = func(path string) error {
		return frontInitAction.FrontInit(sandbox, path)
	}
	actions.FrontPurge = func(path string) error {
		return frontPurgeAction.FrontPurge(sandbox, path)
	}
	actions.AddPage = func(props api.PageProps) error {
		return addPageAction.AddPage(sandbox, props)
	}
	actions.RemovePage = func(path string, name string) error {
		return removePageAction.RemovePage(sandbox, path, name)
	}
	actions.AddDoc = func(props api.DocProps) error {
		return addDocAction.AddDoc(sandbox, props)
	}
	actions.RemoveDoc = func(path string, name string) error {
		return removeDocAction.RemoveDoc(sandbox, path, name)
	}
	actions.AddCliExample = func(path string, name string) error {
		return addCliExampleAction.AddCliExample(sandbox, path, name)
	}
	actions.RemoveCliExample = func(path string, name string) error {
		return removeCliExampleAction.RemoveCliExample(sandbox, path, name)
	}
	actions.AddLibExample = func(path string, name string) error {
		return addLibExampleAction.AddLibExample(sandbox, path, name)
	}
	actions.RemoveLibExample = func(path string, name string) error {
		return removeLibExampleAction.RemoveLibExample(sandbox, path, name)
	}
	actions.ExecTest = func(props api.ExecTestProps) error {
		return execTestsAction.ExecTest(sandbox, props)
	}
	actions.UpdateTest = func(path string, name string) error {
		return updateTestsAction.UpdateTest(sandbox, path, name)
	}
	actions.Interview = func(path string) error {
		return interviewAction.Interview(sandbox, path)
	}

	return actions
}
