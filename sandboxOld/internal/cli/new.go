package cli

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	add_adapter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_adapter"
	add_arg "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_arg"
	add_available "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_available"
	add_body_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_body_field"
	add_cli_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_cli_example"
	add_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_command"
	add_dep "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_dep"
	add_doc "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_doc"
	add_flag "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_flag"
	add_header "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_header"
	add_lib_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_lib_example"
	add_page "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_page"
	add_param "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_param"
	add_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_route"
	add_segment "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_segment"
	build "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/build"
	cli_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli_init"
	cli_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli_purge"
	compile "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/compile"
	deps_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps_init"
	deps_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps_purge"
	exec_test "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/exec_test"
	front_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/front_init"
	front_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/front_purge"
	help "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/help"
	list_adapters "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/list_adapters"
	list_deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/list_deps"
	local_install "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/local_install"
	publish "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/publish"
	remove_adapter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_adapter"
	remove_arg "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_arg"
	remove_available "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_available"
	remove_body_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_body_field"
	remove_cli_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_cli_example"
	remove_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_command"
	remove_dep "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_dep"
	remove_doc "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_doc"
	remove_flag "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_flag"
	remove_header "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_header"
	remove_lib_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_lib_example"
	remove_page "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_page"
	remove_param "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_param"
	remove_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_route"
	remove_segment "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_segment"
	server_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server_init"
	server_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server_purge"
	set_adapter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_adapter"
	set_body "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_body"
	set_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_command"
	set_dep "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_dep"
	set_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_route"
	start "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/start"
	update_test "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/update_test"
	verify "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/verify"
	version "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/version"
)

// NewCli builds the cli surface of the sandbox: Commands, one entry per
// sandbox/internal/commands/<name>/ built by that package's generated
// NewCommand, and CliMain, the dispatch that reads a command line against
// them. Generated by `agnos build` — do not edit by hand.
func NewCli(sandbox *api.Sandbox) api.Cli {
	cli := api.Cli{}

	cli.Commands = []api.Command{
		add_adapter.NewCommand(sandbox),
		add_arg.NewCommand(sandbox),
		add_available.NewCommand(sandbox),
		add_body_field.NewCommand(sandbox),
		add_cli_example.NewCommand(sandbox),
		add_command.NewCommand(sandbox),
		add_dep.NewCommand(sandbox),
		add_doc.NewCommand(sandbox),
		add_flag.NewCommand(sandbox),
		add_header.NewCommand(sandbox),
		add_lib_example.NewCommand(sandbox),
		add_page.NewCommand(sandbox),
		add_param.NewCommand(sandbox),
		add_route.NewCommand(sandbox),
		add_segment.NewCommand(sandbox),
		build.NewCommand(sandbox),
		cli_init.NewCommand(sandbox),
		cli_purge.NewCommand(sandbox),
		compile.NewCommand(sandbox),
		deps_init.NewCommand(sandbox),
		deps_purge.NewCommand(sandbox),
		exec_test.NewCommand(sandbox),
		front_init.NewCommand(sandbox),
		front_purge.NewCommand(sandbox),
		help.NewCommand(sandbox),
		list_adapters.NewCommand(sandbox),
		list_deps.NewCommand(sandbox),
		local_install.NewCommand(sandbox),
		publish.NewCommand(sandbox),
		remove_adapter.NewCommand(sandbox),
		remove_arg.NewCommand(sandbox),
		remove_available.NewCommand(sandbox),
		remove_body_field.NewCommand(sandbox),
		remove_cli_example.NewCommand(sandbox),
		remove_command.NewCommand(sandbox),
		remove_dep.NewCommand(sandbox),
		remove_doc.NewCommand(sandbox),
		remove_flag.NewCommand(sandbox),
		remove_header.NewCommand(sandbox),
		remove_lib_example.NewCommand(sandbox),
		remove_page.NewCommand(sandbox),
		remove_param.NewCommand(sandbox),
		remove_route.NewCommand(sandbox),
		remove_segment.NewCommand(sandbox),
		server_init.NewCommand(sandbox),
		server_purge.NewCommand(sandbox),
		set_adapter.NewCommand(sandbox),
		set_body.NewCommand(sandbox),
		set_command.NewCommand(sandbox),
		set_dep.NewCommand(sandbox),
		set_route.NewCommand(sandbox),
		start.NewCommand(sandbox),
		update_test.NewCommand(sandbox),
		verify.NewCommand(sandbox),
		version.NewCommand(sandbox),
	}

	cli.CliMain = func(args []string) int {
		return CliMain(sandbox, args)
	}

	return cli
}
