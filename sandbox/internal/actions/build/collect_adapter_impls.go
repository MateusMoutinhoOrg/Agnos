package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// CollectAdapterImpls returns one entry per adapters/impls sub-contract
// directory, for the {{range .AdapterImpls}} loop in
// adapters/bindings/standard/new.go.
func CollectAdapterImpls(sandbox *api.Sandbox, io *stagedfs.StagedFS) []map[string]string {
	return collectLibDirs(sandbox, io, "adapters/impls")
}
