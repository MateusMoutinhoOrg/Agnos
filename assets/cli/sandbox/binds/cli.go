package binds

import (
	{{.Module}}/all/sandbox/api"
	"{{.Module}}/all/sandbox/deps"
	"{{.Module}}/sandbox/internal/cli"
)

func CliBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Cli.Commands = append(sandbox.Cli.Commands) //these needs to be constructed

	sandbox.Cli.CliMain = func(args []string) int {
		return cli.CliMain(deps, sandbox, args)
	}
}
