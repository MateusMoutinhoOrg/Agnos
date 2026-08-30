package coremethods

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"

func NewCoreMethods(sandbox *sandbox.SandBox) *CoreMethods {
	return &CoreMethods{
		Sandbox: sandbox,
	}
}
