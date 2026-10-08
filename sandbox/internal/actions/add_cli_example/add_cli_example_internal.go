package add_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddCliExampleInternal writes the one file a cli example holds. The stub it
// renders runs and exits 0 as it stands, so the first run-examples after this
// records a golden instead of reporting a failure. It refuses to overwrite an
// existing example, and refuses outright in a project with no cli: there is no
// binary for an example.sh to type.
func AddCliExampleInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) error {
	if err := utils.RequireExtension(sandbox, io, utils.ExtensionExample); err != nil {
		return err
	}

	if err := utils.ValidateExampleName(sandbox, name); err != nil {
		return err
	}

	has_cli, err := utils.ExtensionEnabled(sandbox, io, utils.ExtensionCli)
	if err != nil {
		return err
	}
	if !has_cli {
		return sandbox.Deps.StdDeps.Errorf("add-cli-example: this project has no cli (%s is off; add one with cli-init)", utils.ExtensionCli)
	}

	dir := utils.ExampleDir(utils.ExampleCliSide, name)
	if io.IsDir(dir) {
		return sandbox.Deps.StdDeps.Errorf("example %s already exists", dir)
	}

	project_conf, err := utils.LoadProjectConf(sandbox, io)
	if err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("add-cli-example creating %s \n", dir)

	vars := map[string]interface{}{
		"ExampleName":   name,
		"GeneratorName": sandbox.Deps.StringsDeps.ToLower(sandbox.Config.ProjectName),
		"ProjectName":   project_conf.ProjectName,
	}
	return utils.RenderTemplateToDest(sandbox, io, "templates/example_cli.sh", vars, dir+"/"+utils.ExampleCliFile)
}
