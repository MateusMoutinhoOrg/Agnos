package sandbox

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"

func Config(sandbox *sandbox.SandBox) {
	sandbox.Config.Version = "v0.0.12"
	sandbox.Config.ProjectName = "Agnos"

}
