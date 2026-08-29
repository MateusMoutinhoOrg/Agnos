package sandbox

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"

func Config(sandbox *lib.SandBox) {
	sandbox.Config.Version = "v0.0.12"
	sandbox.Config.ProjectName = "Agnos"

}
