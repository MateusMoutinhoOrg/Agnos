package deps_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func DepsInitInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Log("deps-init started with path %s \n", path)

	io.CreateDir("sandbox/deps")
	io.CreateDir("adapters")

	return nil
}
