package memory

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func SetMemoryMethods(sandbox *api.SandBox) {
	sandbox.NewMemory = NewMemoryFactory(sandbox)

}
