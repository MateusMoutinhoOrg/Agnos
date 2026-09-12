package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	cli "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/cli"
)

// CliBind wires the generated cli.CliMain (built by `agnos build` from every
// command's entries.yaml) onto the sandbox.
func CliBind(sandbox *api.Sandbox) {
	sandbox.Cli.CliMain = func(args []string) int {
		return cli.CliMain(sandbox, args)
	}
}
