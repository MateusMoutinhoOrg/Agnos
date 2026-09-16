package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CliInitInternal turns the cli mechanic on in the project's declaration. It
// writes no asset itself: the follow-up build renders every group the
// declaration turns on, and sandbox-cli is one of them from here.
func CliInitInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("cli-init started with path %s \n", path)

	return utils.SetExtension(sandbox, io, utils.ExtensionSandboxCli, true)
}
