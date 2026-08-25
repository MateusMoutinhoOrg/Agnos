package memory

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/memory"
)

func NewMemoryFactory(sandbox *api.SandBox) func(path string) Memory {
	return func(path string) Memory {
		memory := memory.Memory{}

		return memory
	}
}

func NewMemoryPakage(sandbox *api.SandBox) {
	sandbox.NewMemory = NewMemoryFactory(sandbox)
}
