package enable_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func EnableDeps(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Printf("enable-deps started with path %s \n", path)

	io.CreateDir("sandbox/deps")
	io.CreateDir("adapters")

	return buildAction.Build(deps, io, path)
}
