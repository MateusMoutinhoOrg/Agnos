package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectAdapterLibs returns one entry per adapters/libs sub-contract
// directory, for the {{range .AdapterLibs}} loop in
// adapters/availables/standard/new.go.
func CollectAdapterLibs(deps *deps.Deps, io *smartio.SmartIO) []map[string]string {
	return collectLibDirs(deps, io, "adapters/libs")
}
