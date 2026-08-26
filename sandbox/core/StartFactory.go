package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/userconfig"
)

func StartFactory(sandbox *api.SandBox) func(path string, project_name string, module *string, force bool) error {

	return func(path string, project_name string, module *string, force bool) error {

		configDir := path + "/" + sandbox.ProjectName + "Config"
		if !force {
			if sandbox.Deps.IoLib.IsDir(configDir) {
				return sandbox.Deps.Errorf("the %s folder already exists", sandbox.ProjectName)
			}
		}

		project_specs_dest := configDir + "/project.yaml"
		project_conf, err := userconfig.NewProjectConf(sandbox, project_specs_dest)
		if err != nil {
			return err
		}

		project_conf.Name = project_name
		project_conf.Description = "Place the Project Description in " + project_specs_dest

		if module != nil {

			module_path := path + "go.mod"

			if !force {
				if sandbox.Deps.IoLib.IsFile(module_path) {
					return sandbox.Deps.Errorf("the %s file already exists", "go.mod")
				}
			}
			module_conf, err := userconfig.NewModuleConf(sandbox, module_path)
			if err != nil {
				return err
			}
			module_conf.Module = *module
			err = module_conf.Persist()
			if err != nil {
				return err
			}
		}

		err = project_conf.Persist()
		if err != nil {
			return err
		}

		sandbox.Deps.Printf("started with path %s \n", path)
		return nil
	}
}
