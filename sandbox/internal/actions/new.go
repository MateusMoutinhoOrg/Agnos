package actions

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_arg"
	addBindingAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_binding"
	addBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_body_field"
	addCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_cli_example"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
	addDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_database"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	addDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_doc"
	addFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_flag"
	addLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_lib_example"
	addPageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_page"
	addParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_parameter"
	addPathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_path"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
	addTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table"
	addTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table_field"
	backofficeInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/backoffice_init"
	backofficePurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/backoffice_purge"
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
	explainCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/explain_command"
	explainRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/explain_route"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
	frontPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_purge"
	importBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/import_body"
	interviewAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/interview"
	listAdaptersAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_adapters"
	listCommandsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_commands"
	listDepsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_deps"
	listExtensionsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_extensions"
	listRoutesAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_routes"
	rebalanceCommandsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rebalance_commands"
	rebalanceRoutesAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rebalance_routes"
	removeAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_adapter"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_arg"
	removeBindingAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_binding"
	removeBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_body_field"
	removeCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_cli_example"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_command"
	removeDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_database"
	removeDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_dep"
	removeDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_doc"
	removeFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_flag"
	removeLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_lib_example"
	removePageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_page"
	removeParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_parameter"
	removePathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_path"
	removeRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_route"
	removeTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table"
	removeTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table_field"
	renameCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rename_command"
	renameRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rename_route"
	runExamplesAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/run_examples"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	serverPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_purge"
	setAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_adapter"
	setArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_arg"
	setBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body"
	setBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body_field"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_command"
	setDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_dep"
	setFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_flag"
	setParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_parameter"
	setPathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_path"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
	setTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_table_field"
	showCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_command"
	showDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_database"
	showRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_route"
	startAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/start"
	updateExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/update_example"
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
	actions.Verify = func(props api.VerifyProps) error {
		return verifyAction.Verify(sandbox, props)
	}
	actions.Start = func(props api.StartProps) error {
		return startAction.Start(sandbox, props)
	}
	actions.EnableExtension = func(props api.EnableExtensionProps) error {
		return enableExtensionAction.EnableExtension(sandbox, props)
	}
	actions.DisableExtension = func(props api.DisableExtensionProps) error {
		return disableExtensionAction.DisableExtension(sandbox, props)
	}
	actions.ListExtensions = func(props api.ListExtensionsProps) ([]api.ExtensionInfo, error) {
		return listExtensionsAction.ListExtensions(sandbox, props)
	}
	actions.DepsInit = func(props api.DepsInitProps) error {
		return depsInitAction.DepsInit(sandbox, props)
	}
	actions.DepsPurge = func(props api.DepsPurgeProps) error {
		return depsPurgeAction.DepsPurge(sandbox, props)
	}
	actions.AddDep = func(props api.AddDepProps) error {
		return addDepAction.AddDep(sandbox, props)
	}
	actions.RemoveDep = func(props api.RemoveDepProps) error {
		return removeDepAction.RemoveDep(sandbox, props)
	}
	actions.ListDeps = func(props api.ListDepsProps) ([]api.DepInfo, error) {
		return listDepsAction.ListDeps(sandbox, props)
	}
	actions.SetDep = func(props api.SetDepProps) error {
		return setDepAction.SetDep(sandbox, props)
	}
	actions.AddAdapter = func(props api.AddAdapterProps) error {
		return addAdapterAction.AddAdapter(sandbox, props)
	}
	actions.RemoveAdapter = func(props api.RemoveAdapterProps) error {
		return removeAdapterAction.RemoveAdapter(sandbox, props)
	}
	actions.SetAdapter = func(props api.SetAdapterProps) error {
		return setAdapterAction.SetAdapter(sandbox, props)
	}
	actions.ListAdapters = func(props api.ListAdaptersProps) ([]api.AdapterInfo, error) {
		return listAdaptersAction.ListAdapters(sandbox, props)
	}
	actions.AddBinding = func(props api.AddBindingProps) error {
		return addBindingAction.AddBinding(sandbox, props)
	}
	actions.RemoveBinding = func(props api.RemoveBindingProps) error {
		return removeBindingAction.RemoveBinding(sandbox, props)
	}
	actions.CliInit = func(props api.CliInitProps) error {
		return cliInitAction.CliInit(sandbox, props)
	}
	actions.CliPurge = func(props api.CliPurgeProps) error {
		return cliPurgeAction.CliPurge(sandbox, props)
	}
	actions.AddCommand = func(props api.AddCommandProps) error {
		return addCommandAction.AddCommand(sandbox, props)
	}
	actions.RemoveCommand = func(props api.RemoveCommandProps) error {
		return removeCommandAction.RemoveCommand(sandbox, props)
	}
	actions.SetCommand = func(props api.SetCommandProps) error {
		return setCommandAction.SetCommand(sandbox, props)
	}
	actions.RenameCommand = func(props api.RenameCommandProps) error {
		return renameCommandAction.RenameCommand(sandbox, props)
	}
	actions.RebalanceCommands = func(props api.RebalanceCommandsProps) error {
		return rebalanceCommandsAction.RebalanceCommands(sandbox, props)
	}
	actions.ListCommands = func(props api.ListCommandsProps) ([]string, error) {
		return listCommandsAction.ListCommands(sandbox, props)
	}
	actions.ShowCommand = func(props api.ShowCommandProps) ([]string, error) {
		return showCommandAction.ShowCommand(sandbox, props)
	}
	actions.ExplainCommand = func(props api.ExplainCommandProps) ([]string, error) {
		return explainCommandAction.ExplainCommand(sandbox, props)
	}
	actions.AddFlag = func(props api.AddFlagProps) error {
		return addFlagAction.AddFlag(sandbox, props)
	}
	actions.SetFlag = func(props api.SetFlagProps) error {
		return setFlagAction.SetFlag(sandbox, props)
	}
	actions.RemoveFlag = func(props api.RemoveFlagProps) error {
		return removeFlagAction.RemoveFlag(sandbox, props)
	}
	actions.AddArg = func(props api.AddArgProps) error {
		return addArgAction.AddArg(sandbox, props)
	}
	actions.SetArg = func(props api.SetArgProps) error {
		return setArgAction.SetArg(sandbox, props)
	}
	actions.RemoveArg = func(props api.RemoveArgProps) error {
		return removeArgAction.RemoveArg(sandbox, props)
	}
	actions.ServerInit = func(props api.ServerInitProps) error {
		return serverInitAction.ServerInit(sandbox, props)
	}
	actions.ServerPurge = func(props api.ServerPurgeProps) error {
		return serverPurgeAction.ServerPurge(sandbox, props)
	}
	actions.AddRoute = func(props api.AddRouteProps) error {
		return addRouteAction.AddRoute(sandbox, props)
	}
	actions.RemoveRoute = func(props api.RemoveRouteProps) error {
		return removeRouteAction.RemoveRoute(sandbox, props)
	}
	actions.SetRoute = func(props api.SetRouteProps) error {
		return setRouteAction.SetRoute(sandbox, props)
	}
	actions.AddPath = func(props api.AddPathProps) error {
		return addPathAction.AddPath(sandbox, props)
	}
	actions.SetPath = func(props api.SetPathProps) error {
		return setPathAction.SetPath(sandbox, props)
	}
	actions.RemovePath = func(props api.RemovePathProps) error {
		return removePathAction.RemovePath(sandbox, props)
	}
	actions.AddParameter = func(props api.AddParameterProps) error {
		return addParameterAction.AddParameter(sandbox, props)
	}
	actions.SetParameter = func(props api.SetParameterProps) error {
		return setParameterAction.SetParameter(sandbox, props)
	}
	actions.RemoveParameter = func(props api.RemoveParameterProps) error {
		return removeParameterAction.RemoveParameter(sandbox, props)
	}
	actions.SetBody = func(props api.SetBodyProps) error {
		return setBodyAction.SetBody(sandbox, props)
	}
	actions.AddBodyField = func(props api.AddBodyFieldProps) error {
		return addBodyFieldAction.AddBodyField(sandbox, props)
	}
	actions.RemoveBodyField = func(props api.RemoveBodyFieldProps) error {
		return removeBodyFieldAction.RemoveBodyField(sandbox, props)
	}
	actions.SetBodyField = func(props api.SetBodyFieldProps) error {
		return setBodyFieldAction.SetBodyField(sandbox, props)
	}
	actions.ImportBody = func(props api.ImportBodyProps) error {
		return importBodyAction.ImportBody(sandbox, props)
	}
	actions.ShowRoute = func(props api.ShowRouteProps) ([]string, error) {
		return showRouteAction.ShowRoute(sandbox, props)
	}
	actions.ListRoutes = func(props api.ListRoutesProps) ([]string, error) {
		return listRoutesAction.ListRoutes(sandbox, props)
	}
	actions.ExplainRoute = func(props api.ExplainRouteProps) ([]string, error) {
		return explainRouteAction.ExplainRoute(sandbox, props)
	}
	actions.RenameRoute = func(props api.RenameRouteProps) error {
		return renameRouteAction.RenameRoute(sandbox, props)
	}
	actions.RebalanceRoutes = func(props api.RebalanceRoutesProps) error {
		return rebalanceRoutesAction.RebalanceRoutes(sandbox, props)
	}
	actions.DatabaseInit = func(props api.DatabaseInitProps) error {
		return databaseInitAction.DatabaseInit(sandbox, props)
	}
	actions.DatabasePurge = func(props api.DatabasePurgeProps) error {
		return databasePurgeAction.DatabasePurge(sandbox, props)
	}
	actions.AddDatabase = func(props api.AddDatabaseProps) error {
		return addDatabaseAction.AddDatabase(sandbox, props)
	}
	actions.RemoveDatabase = func(props api.RemoveDatabaseProps) error {
		return removeDatabaseAction.RemoveDatabase(sandbox, props)
	}
	actions.AddTable = func(props api.AddTableProps) error {
		return addTableAction.AddTable(sandbox, props)
	}
	actions.RemoveTable = func(props api.RemoveTableProps) error {
		return removeTableAction.RemoveTable(sandbox, props)
	}
	actions.AddTableField = func(props api.AddTableFieldProps) error {
		return addTableFieldAction.AddTableField(sandbox, props)
	}
	actions.SetTableField = func(props api.SetTableFieldProps) error {
		return setTableFieldAction.SetTableField(sandbox, props)
	}
	actions.RemoveTableField = func(props api.RemoveTableFieldProps) error {
		return removeTableFieldAction.RemoveTableField(sandbox, props)
	}
	actions.ShowDatabase = func(props api.ShowDatabaseProps) ([]string, error) {
		return showDatabaseAction.ShowDatabase(sandbox, props)
	}
	actions.FrontInit = func(props api.FrontInitProps) error {
		return frontInitAction.FrontInit(sandbox, props)
	}
	actions.FrontPurge = func(props api.FrontPurgeProps) error {
		return frontPurgeAction.FrontPurge(sandbox, props)
	}
	actions.BackofficeInit = func(props api.BackofficeInitProps) error {
		return backofficeInitAction.BackofficeInit(sandbox, props)
	}
	actions.BackofficePurge = func(props api.BackofficePurgeProps) error {
		return backofficePurgeAction.BackofficePurge(sandbox, props)
	}
	actions.AddPage = func(props api.AddPageProps) error {
		return addPageAction.AddPage(sandbox, props)
	}
	actions.RemovePage = func(props api.RemovePageProps) error {
		return removePageAction.RemovePage(sandbox, props)
	}
	actions.AddDoc = func(props api.AddDocProps) error {
		return addDocAction.AddDoc(sandbox, props)
	}
	actions.RemoveDoc = func(props api.RemoveDocProps) error {
		return removeDocAction.RemoveDoc(sandbox, props)
	}
	actions.AddCliExample = func(props api.AddCliExampleProps) error {
		return addCliExampleAction.AddCliExample(sandbox, props)
	}
	actions.RemoveCliExample = func(props api.RemoveCliExampleProps) error {
		return removeCliExampleAction.RemoveCliExample(sandbox, props)
	}
	actions.AddLibExample = func(props api.AddLibExampleProps) error {
		return addLibExampleAction.AddLibExample(sandbox, props)
	}
	actions.RemoveLibExample = func(props api.RemoveLibExampleProps) error {
		return removeLibExampleAction.RemoveLibExample(sandbox, props)
	}
	actions.RunExamples = func(props api.RunExamplesProps) error {
		return runExamplesAction.RunExamples(sandbox, props)
	}
	actions.UpdateExample = func(props api.UpdateExampleProps) error {
		return updateExampleAction.UpdateExample(sandbox, props)
	}
	actions.Interview = func(props api.InterviewProps) error {
		return interviewAction.Interview(sandbox, props)
	}

	return actions
}
