package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// CliInitInternal renders every embedded asset under assets/cli into the
// target project at the path it holds inside that group, using the same
// Module variable the build step derives from go.mod.
func CliInitInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Printf("cli-init started with path %s \n", path)

	gomod, err := io.ReadFile("go.mod")
	if err != nil {
		return err
	}
	module_conf, err := moduleconf.New(deps, string(gomod))
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module": module_conf.Module,
	}

	return utils.RenderGroup(deps, io, "cli", vars)
}
