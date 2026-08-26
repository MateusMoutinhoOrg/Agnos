package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func SetCoreMethods(sandbox *api.SandBox) {
	sandbox.Start = StartFactory(sandbox)
	sandbox.Build = BuildFactory(sandbox)
}
