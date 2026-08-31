package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func ExportItens(deps deps.Deps, sandbox *api.Sandbox) {
	sandbox.Cli.CliMain = func(args []string) int {
		return CliMain(deps, args)
	}
}
