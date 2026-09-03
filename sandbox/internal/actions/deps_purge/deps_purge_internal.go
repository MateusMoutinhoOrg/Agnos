package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func DepsPurgeInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Log("deps-purge started with path %s \n", path)

	io.RemoveDir("sandbox/deps")
	io.RemoveDir("adapters")

	return nil
}
