package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// CollectAdapterLibs returns one entry per adapters/libs sub-contract
// directory, for the {{range .AdapterLibs}} loop in
// adapters/availables/standard/new.go.
func CollectAdapterLibs(io *smartio.SmartIO) []map[string]string {
	return collectLibDirs(io, "adapters/libs")
}
