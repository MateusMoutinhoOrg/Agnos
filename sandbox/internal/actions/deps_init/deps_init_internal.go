package deps_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/bindingconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

func DepsInitInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("deps-init started with path %s \n", path)

	io.CreateDir("sandbox/deps")
	io.CreateDir("adapters")

	if err := utils.SetExtension(sandbox, io, utils.ExtensionDeps, true); err != nil {
		return err
	}

	// The standard binding starts empty and declared: from here on which
	// adapter binds is read from the declaration, never from a listing of
	// adapters/impls, so the same contract may later have two implementations.
	if io.IsFile(utils.BindingConfPath(utils.StandardBinding)) {
		return nil
	}

	return io.WriteFile(utils.BindingConfPath(utils.StandardBinding),
		[]byte(bindingconf.NewEmpty(sandbox).Render()))
}
