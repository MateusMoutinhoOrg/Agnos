package remove_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func RemoveDeps(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Printf("remove-deps started with path %s \n", path)

	io.RemoveDir("sandbox/deps")
	io.RemoveDir("adapters")

	return nil
}
