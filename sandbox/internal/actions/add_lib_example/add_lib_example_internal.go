package add_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddLibExampleInternal writes the one file a lib example holds. The stub it
// renders runs and exits 0 as it stands, so the first exec-test after this
// records a golden instead of reporting a failure. It refuses to overwrite an
// existing example.
func AddLibExampleInternal(deps *deps.Deps, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateExampleName(deps, name); err != nil {
		return err
	}

	dir := utils.ExampleDir(utils.ExampleLibSide, name)
	if io.IsDir(dir) {
		return deps.Std.Errorf("example %s already exists", dir)
	}

	module_conf, err := utils.LoadModuleConf(deps, io)
	if err != nil {
		return err
	}

	deps.Std.Log("add-lib-example creating %s \n", dir)

	vars := map[string]interface{}{
		"Name":    name,
		"Module":  module_conf.Module,
		"HasDeps": io.IsDir("sandbox/deps"),
	}
	return utils.RenderTemplateToDest(deps, io, "templates/example_lib.go", vars, dir+"/"+utils.ExampleLibFile)
}
