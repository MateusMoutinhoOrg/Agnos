package cli

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"

func NewCli(sandbox *api.SandBox) {

	sandbox.CliMain = CliMainFactory(sandbox)
}
