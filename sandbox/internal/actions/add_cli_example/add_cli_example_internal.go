package add_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddCliExampleInternal writes the one file a cli example holds. The stub it
// renders runs and exits 0 as it stands, so the first exec-test after this
// records a golden instead of reporting a failure. It refuses to overwrite an
// existing example, and refuses outright in a project with no cli: there is no
// binary for an example.sh to type.
func AddCliExampleInternal(deps *deps.Deps, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateExampleName(deps, name); err != nil {
		return err
	}

	if !io.IsDir("sandbox/internal/cli") {
		return deps.Std.Errorf("add-cli-example: this project has no cli (sandbox/internal/cli is missing; add one with cli-init)")
	}

	dir := utils.ExampleDir(utils.ExampleCliSide, name)
	if io.IsDir(dir) {
		return deps.Std.Errorf("example %s already exists", dir)
	}

	project_conf, err := utils.LoadProjectConf(deps, io)
	if err != nil {
		return err
	}

	deps.Std.Log("add-cli-example creating %s \n", dir)

	vars := map[string]interface{}{
		"Name":        name,
		"ProjectName": project_conf.Name,
	}
	return utils.RenderTemplateToDest(deps, io, "templates/example_cli.sh", vars, dir+"/"+utils.ExampleCliFile)
}
