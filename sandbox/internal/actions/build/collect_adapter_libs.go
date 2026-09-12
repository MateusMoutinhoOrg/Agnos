package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectAdapterLibs returns one entry per adapters/libs sub-contract
// directory, for the {{range .AdapterLibs}} loop in
// adapters/availables/standard/new.go.
func CollectAdapterLibs(sandbox *api.Sandbox, io *smartio.SmartIO) []map[string]string {
	return collectLibDirs(sandbox, io, "adapters/libs")
}
