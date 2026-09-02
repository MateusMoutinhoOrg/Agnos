package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

func BuildInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Printf("build started with path %s \n", path)

	//Creating the basic dir struct
	io.CreateDir("sandbox/api")
	io.CreateDir("sandbox/internal")

	gomod, err := io.ReadFile(path + "/go.mod")
	if err != nil {
		return err
	}
	module_conf, err := moduleconf.New(deps, string(gomod))
	if err != nil {
		return err
	}

	hasDeps := io.IsDir("sandbox/deps")
	hasCli := io.IsDir("sandbox/internal/cli")

	vars := map[string]interface{}{
		"Module":       module_conf.Module,
		"HasDeps":      hasDeps,
		"HasCli":       hasCli,
		"Binds":        CollectBinds(io),
		"Constructors": CollectConstructors(io),
		"DepsLibs":     CollectDepsLibs(io),
		"AdapterLibs":  CollectAdapterLibs(io),
		"Commands":     CollectCommands(io),
	}

	if err := utils.RenderGroup(deps, io, "all", vars); err != nil {
		return err
	}

	if hasDeps {
		if err := utils.RenderGroup(deps, io, "deps", vars); err != nil {
			return err
		}
	}

	if hasCli {
		if err := utils.RenderGroup(deps, io, "cli", vars); err != nil {
			return err
		}
	}

	deps.Std.Printf("successfully rendered template")
	return nil
}
