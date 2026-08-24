package cli

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"

func SetCliProps(sandbox *api.SandBox) {
	sandbox.CliMain = CliMainFactory(sandbox)
}
