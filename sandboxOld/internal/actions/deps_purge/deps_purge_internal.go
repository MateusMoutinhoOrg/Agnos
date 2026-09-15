package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func DepsPurgeInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("deps-purge started with path %s \n", path)

	io.RemoveDir("sandbox/deps")
	io.RemoveDir("adapters")

	return nil
}
