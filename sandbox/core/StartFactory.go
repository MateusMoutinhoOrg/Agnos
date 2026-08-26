package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/userconfig"
)

func StartFactory(sandbox *api.SandBox) func(props api.StartProps) error {

	return func(props api.StartProps) error {

		configDir := props.Path + "/" + sandbox.ProjectName + "Config"
		if !props.Force {
			if sandbox.Deps.IoLib.IsDir(configDir) {
				return sandbox.Deps.Errorf("the %s folder already exists", sandbox.ProjectName)
			}
		}

		project_specs_dest := configDir + "/project.yaml"
		project_conf, err := userconfig.NewProjectConf(sandbox, project_specs_dest)
		if err != nil {
			return err
		}

		project_conf.Name = props.ProjectName
		project_conf.Description = "Place the Project Description in " + project_specs_dest

		if props.Module != nil {

			module_path := props.Path + "/go.mod"

			if !props.Force {
				if sandbox.Deps.IoLib.IsFile(module_path) {
					return sandbox.Deps.Errorf("the %s file already exists", "go.mod")
				}
			}
			module_conf, err := userconfig.NewModuleConf(sandbox, module_path)
			if err != nil {
				return err
			}
			module_conf.Module = *props.Module
			err = module_conf.Persist()
			if err != nil {
				return err
			}
		}

		err = project_conf.Persist()
		if err != nil {
			return err
		}

		sandbox.Deps.Printf("started with path %s \n", props.Path)
		return nil
	}
}
