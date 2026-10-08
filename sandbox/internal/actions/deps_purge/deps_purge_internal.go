package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

func DepsPurgeInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("deps-purge started with path %s \n", path)

	io.RemoveDir("sandbox/deps")
	io.RemoveDir("adapters")

	return utils.SetExtension(sandbox, io, utils.ExtensionDeps, false)
}
