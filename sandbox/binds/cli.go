package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	cli "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/cli"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/build"
	dep_install "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/dep_install"
	dep_list "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/dep_list"
	dep_remove "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/dep_remove"
	deps_init "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/deps_init"
	deps_purge "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/deps_purge"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/help"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/start"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/version"
)

func CliBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Cli.Commands = append(sandbox.Cli.Commands,
		build.NewCommand(deps, sandbox),
		start.NewCommand(deps, sandbox),
		deps_init.NewCommand(deps, sandbox),
		deps_purge.NewCommand(deps, sandbox),
		dep_install.NewCommand(deps, sandbox),
		dep_remove.NewCommand(deps, sandbox),
		dep_list.NewCommand(deps, sandbox),
		version.NewCommand(deps, sandbox),
		help.NewCommand(deps, sandbox),
	)

	sandbox.Cli.CliMain = func(args []string) int {
		return cli.CliMain(deps, sandbox, args)
	}
}
