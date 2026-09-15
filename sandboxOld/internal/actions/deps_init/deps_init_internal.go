package deps_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/availableconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

func DepsInitInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("deps-init started with path %s \n", path)

	io.CreateDir("sandbox/deps")
	io.CreateDir("adapters")

	// The standard available starts empty and declared: from here on which
	// adapter binds is read from the declaration, never from a listing of
	// adapters/libs, so the same contract may later have two implementations.
	if io.IsFile(utils.AvailableConfPath(utils.StandardAvailable)) {
		return nil
	}

	return io.WriteFileOverwrite(utils.AvailableConfPath(utils.StandardAvailable),
		[]byte(availableconf.NewEmpty(sandbox).Render()))
}
