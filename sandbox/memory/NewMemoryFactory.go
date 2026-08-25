package memory

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewMemoryFactory(sandbox *api.SandBox) func(path string) api.Memory {
	return func(path string) api.Memory {

		memory := api.Memory{}

		return memory
	}
}
