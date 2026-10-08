package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CliInitInternal turns the cli mechanic on in the project's declaration. It
// writes no asset itself: the follow-up build renders every group the
// declaration turns on, and cli is one of them from here.
func CliInitInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("cli-init started with path %s \n", path)

	return utils.SetExtension(sandbox, io, utils.ExtensionCli, true)
}
