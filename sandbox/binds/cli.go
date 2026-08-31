package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	cli "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/cli"
	
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/help"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/start"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/commands/version"
)

func CliBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Cli.Commands = append(sandbox.Cli.Commands,
		build.NewCommand(deps, sandbox),
		version.NewCommand(deps, sandbox),
		help.NewCommand(deps, sandbox),
		start.NewCommand(deps, sandbox),
	)

	sandbox.Cli.CliMain = func(args []string) int {
		return cli.CliMain(deps, sandbox, args)
	}
}
