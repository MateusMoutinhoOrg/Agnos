package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	cli "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/cli"
)

func CliBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Cli.CliMain = func(args []string) int {
		return cli.CliMain(deps, args)
	}
}
