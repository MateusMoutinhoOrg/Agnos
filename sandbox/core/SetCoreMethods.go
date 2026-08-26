package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/start"
)

func SetCoreMethods(sandbox *api.SandBox) {
	sandbox.Start = start.StartFactory(sandbox)
	sandbox.Build = build.BuildFactory(sandbox)
}
