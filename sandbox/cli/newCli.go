package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/version"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewCli(sandbox *api.SandBox) {

	sandbox.Commands = []api.CliCommand{
		version.NewCommand(sandbox),
	}

	sandbox.CliMain = CliMainFactory(sandbox)
}
