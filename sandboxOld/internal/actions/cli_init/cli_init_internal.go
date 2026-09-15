package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CliInitInternal renders every embedded asset under assets/cli into the
// target project at the path it holds inside that group, using the same
// Module variable the build step derives from go.mod.
func CliInitInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("cli-init started with path %s \n", path)

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module": module_conf.Module,
	}

	return utils.RenderGroup(sandbox, io, "cli", vars)
}
