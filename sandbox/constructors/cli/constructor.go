package cli

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	cli "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/cli"
)

func Constructor(sandbox *api.Sandbox) {

	sandbox.Cli = cli.NewCli(sandbox)
}
