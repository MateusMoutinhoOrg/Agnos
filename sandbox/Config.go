package sandbox

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"

func Config(sandbox *api.SandBox) {
	sandbox.Version = "v0.0.10"
	sandbox.ProjectName = "Agnos"
}
