package enable_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func EnableDepsInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Printf("enable-deps started with path %s \n", path)

	io.CreateDir("sandbox/deps")
	io.CreateDir("adapters")

	return nil
}
