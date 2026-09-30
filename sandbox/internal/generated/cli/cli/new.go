package cli

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/cli/errors"
	add_arg "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/add_arg"
	add_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/add_command"
	add_flag "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/add_flag"
	cli_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/cli_init"
	cli_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/cli_purge"
	explain_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/explain_command"
	list_commands "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/list_commands"
	rebalance_commands "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/rebalance_commands"
	remove_arg "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/remove_arg"
	remove_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/remove_command"
	remove_flag "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/remove_flag"
	rename_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/rename_command"
	set_arg "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/set_arg"
	set_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/set_command"
	set_flag "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/set_flag"
	show_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli/show_command"
	build "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/core/build"
	compile "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/core/compile"
	interview "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/core/interview"
	local_install "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/core/local_install"
	publish "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/core/publish"
	start "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/core/start"
	verify "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/core/verify"
	add_database "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/add_database"
	add_table "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/add_table"
	add_table_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/add_table_field"
	database_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/database_init"
	database_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/database_purge"
	remove_database "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/remove_database"
	remove_table "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/remove_table"
	remove_table_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/remove_table_field"
	set_table_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/set_table_field"
	show_database "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/database/show_database"
	add_adapter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/add_adapter"
	add_available "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/add_available"
	add_dep "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/add_dep"
	deps_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/deps_init"
	deps_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/deps_purge"
	list_adapters "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/list_adapters"
	list_deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/list_deps"
	remove_adapter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/remove_adapter"
	remove_available "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/remove_available"
	remove_dep "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/remove_dep"
	set_adapter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/set_adapter"
	set_dep "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps/set_dep"
	add_doc "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/docs/add_doc"
	remove_doc "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/docs/remove_doc"
	add_cli_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/examples/add_cli_example"
	add_lib_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/examples/add_lib_example"
	exec_test "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/examples/exec_test"
	remove_cli_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/examples/remove_cli_example"
	remove_lib_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/examples/remove_lib_example"
	update_test "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/examples/update_test"
	disable_extension "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/extensions/disable_extension"
	enable_extension "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/extensions/enable_extension"
	list_extensions "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/extensions/list_extensions"
	add_page "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/front/add_page"
	front_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/front/front_init"
	front_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/front/front_purge"
	remove_page "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/front/remove_page"
	help "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/help"
	help_flag "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/help_flag"
	project "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/project"
	add_body_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/add_body_field"
	add_parameter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/add_parameter"
	add_path "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/add_path"
	add_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/add_route"
	explain_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/explain_route"
	import_body "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/import_body"
	list_routes "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/list_routes"
	rebalance_routes "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/rebalance_routes"
	remove_body_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/remove_body_field"
	remove_parameter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/remove_parameter"
	remove_path "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/remove_path"
	remove_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/remove_route"
	rename_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/rename_route"
	server_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/server_init"
	server_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/server_purge"
	set_body "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/set_body"
	set_body_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/set_body_field"
	set_parameter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/set_parameter"
	set_path "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/set_path"
	set_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/set_route"
	show_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server/show_route"
	version "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/version"
)

// NewCli builds the cli surface of the sandbox: Commands, one entry per
// directory under sandbox/internal/commands holding a command.yaml, at any
// depth, built by that package's generated NewCommand, in run order; Fail,
// which hands a failure to the project's own handler for it; and CliMain, the
// dispatch that reads a command line against them. Generated by
// `agnos build` — do not edit by hand.
func NewCli(sandbox *api.Sandbox) api.Cli {
	cli := api.Cli{}

	cli.Commands = []*api.Command{
		help_flag.NewCommand(sandbox),
		project.NewCommand(sandbox),
		add_adapter.NewCommand(sandbox),
		add_arg.NewCommand(sandbox),
		add_available.NewCommand(sandbox),
		add_body_field.NewCommand(sandbox),
		add_cli_example.NewCommand(sandbox),
		add_command.NewCommand(sandbox),
		add_database.NewCommand(sandbox),
		add_dep.NewCommand(sandbox),
		add_doc.NewCommand(sandbox),
		add_flag.NewCommand(sandbox),
		add_lib_example.NewCommand(sandbox),
		add_page.NewCommand(sandbox),
		add_parameter.NewCommand(sandbox),
		add_path.NewCommand(sandbox),
		add_route.NewCommand(sandbox),
		add_table.NewCommand(sandbox),
		add_table_field.NewCommand(sandbox),
		build.NewCommand(sandbox),
		cli_init.NewCommand(sandbox),
		cli_purge.NewCommand(sandbox),
		compile.NewCommand(sandbox),
		database_init.NewCommand(sandbox),
		database_purge.NewCommand(sandbox),
		deps_init.NewCommand(sandbox),
		deps_purge.NewCommand(sandbox),
		disable_extension.NewCommand(sandbox),
		enable_extension.NewCommand(sandbox),
		exec_test.NewCommand(sandbox),
		explain_command.NewCommand(sandbox),
		explain_route.NewCommand(sandbox),
		front_init.NewCommand(sandbox),
		front_purge.NewCommand(sandbox),
		help.NewCommand(sandbox),
		import_body.NewCommand(sandbox),
		interview.NewCommand(sandbox),
		list_adapters.NewCommand(sandbox),
		list_commands.NewCommand(sandbox),
		list_deps.NewCommand(sandbox),
		list_extensions.NewCommand(sandbox),
		list_routes.NewCommand(sandbox),
		local_install.NewCommand(sandbox),
		publish.NewCommand(sandbox),
		rebalance_commands.NewCommand(sandbox),
		rebalance_routes.NewCommand(sandbox),
		remove_adapter.NewCommand(sandbox),
		remove_arg.NewCommand(sandbox),
		remove_available.NewCommand(sandbox),
		remove_body_field.NewCommand(sandbox),
		remove_cli_example.NewCommand(sandbox),
		remove_command.NewCommand(sandbox),
		remove_database.NewCommand(sandbox),
		remove_dep.NewCommand(sandbox),
		remove_doc.NewCommand(sandbox),
		remove_flag.NewCommand(sandbox),
		remove_lib_example.NewCommand(sandbox),
		remove_page.NewCommand(sandbox),
		remove_parameter.NewCommand(sandbox),
		remove_path.NewCommand(sandbox),
		remove_route.NewCommand(sandbox),
		remove_table.NewCommand(sandbox),
		remove_table_field.NewCommand(sandbox),
		rename_command.NewCommand(sandbox),
		rename_route.NewCommand(sandbox),
		server_init.NewCommand(sandbox),
		server_purge.NewCommand(sandbox),
		set_adapter.NewCommand(sandbox),
		set_arg.NewCommand(sandbox),
		set_body.NewCommand(sandbox),
		set_body_field.NewCommand(sandbox),
		set_command.NewCommand(sandbox),
		set_dep.NewCommand(sandbox),
		set_flag.NewCommand(sandbox),
		set_parameter.NewCommand(sandbox),
		set_path.NewCommand(sandbox),
		set_route.NewCommand(sandbox),
		set_table_field.NewCommand(sandbox),
		show_command.NewCommand(sandbox),
		show_database.NewCommand(sandbox),
		show_route.NewCommand(sandbox),
		start.NewCommand(sandbox),
		update_test.NewCommand(sandbox),
		verify.NewCommand(sandbox),
		version.NewCommand(sandbox),
	}

	cli.Fail = func(command *api.Command) error {
		if command.Failure == nil {
			return errors.HandleFailure(sandbox, command, command.Response)
		}
		switch command.Failure.Kind {
		case api.NotFoundFailure:
			return errors.HandleNotFound(sandbox, command, command.Response)
		case api.BadUsageFailure:
			return errors.HandleBadUsage(sandbox, command, command.Response)
		case api.UnknownFlagFailure:
			return errors.HandleUnknownFlag(sandbox, command, command.Response)
		case api.UnexpectedArgFailure:
			return errors.HandleUnexpectedArg(sandbox, command, command.Response)
		}
		return errors.HandleFailure(sandbox, command, command.Response)
	}

	cli.CliMain = func(args []string) int {
		return CliMain(sandbox, args)
	}

	return cli
}
