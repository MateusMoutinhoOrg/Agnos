package add_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddLibExampleInternal writes the one file a lib example holds. The stub it
// renders runs and exits 0 as it stands, so the first run-examples after this
// records a golden instead of reporting a failure. It refuses to overwrite an
// existing example.
func AddLibExampleInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) error {
	if err := utils.RequireExtension(sandbox, io, utils.ExtensionExample); err != nil {
		return err
	}

	if err := utils.ValidateExampleName(sandbox, name); err != nil {
		return err
	}

	dir := utils.ExampleDir(utils.ExampleLibSide, name)
	if io.IsDir(dir) {
		return sandbox.Deps.StdDeps.Errorf("example %s already exists", dir)
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("add-lib-example creating %s \n", dir)

	has_deps, err := utils.ExtensionEnabled(sandbox, io, utils.ExtensionDeps)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"ExampleName": name,
		"Module":      module_conf.Module,
		"HasDeps":     has_deps,
	}
	return utils.RenderTemplateToDest(sandbox, io, "templates/example_lib.go", vars, dir+"/"+utils.ExampleLibFile)
}
