package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/parsables"
)

func createProjectSpecs(sandbox *api.SandBox, props api.StartProps, configDir string) error {

	project_specs_dest := configDir + "/project.yaml"

	var project_conf *parsables.ProjectConf
	var err error

	if sandbox.Deps.IoLib.IsFile(project_specs_dest) {
		content_bytes, fileerror := sandbox.Deps.IoLib.ReadFile(project_specs_dest)
		if fileerror != nil {
			return fileerror
		}
		project_conf, err = parsables.NewProjectConf(sandbox, string(content_bytes))
		if err != nil {
			return err
		}
	} else {
		project_conf = parsables.NewProjectConfEmpty(sandbox)
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

		var module_conf *parsables.ModuleConf

		if sandbox.Deps.IoLib.IsFile(module_path) {
			content_bytes, fileerror := sandbox.Deps.IoLib.ReadFile(module_path)
			if fileerror != nil {
				return fileerror
			}
			module_conf, err = parsables.NewModuleConf(sandbox, string(content_bytes))
			if err != nil {
				return err
			}
		} else {
			module_conf = parsables.NewModuleConfEmpty(sandbox)
		}

		module_conf.Module = *props.Module
		rendered := module_conf.Render()
		err = sandbox.Deps.IoLib.WriteFile(module_path, []byte(rendered))
		if err != nil {
			return err
		}
	}

	rendered := project_conf.Render()
	err = sandbox.Deps.IoLib.WriteFile(project_specs_dest, []byte(rendered))
	if err != nil {
		return err
	}
	return nil
}
