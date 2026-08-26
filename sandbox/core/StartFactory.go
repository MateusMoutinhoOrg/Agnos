package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/userconfig"
)

func StartFactory(sandbox *api.SandBox) func(path string, project_name string, module string) error {

	return func(path string, project_name string, module string) error {

		configDir := path + "/" + sandbox.ProjectName + "Config"
		if sandbox.Deps.IoLib.IsDir(configDir) {
			return sandbox.Deps.Errorf("the %s folder already exists", sandbox.ProjectName)
		}

		project_specs_dest := configDir + "/project.yaml"
		project_conf, err := userconfig.NewProjectConf(sandbox, project_specs_dest)
		if err != nil {
			return err
		}

		project_conf.Name = project_name
		project_conf.Module = module
		project_conf.Description = "Place the Project Description in " + project_specs_dest
		err = project_conf.Persist()
		if err != nil {
			return err
		}


		sandbox.Deps.Printf("started with path %s \n", path)
		return nil
	}
}
