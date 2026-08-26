package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func StartFactory(sandbox *api.SandBox) func(path string, project_name string, module string) error {

	return func(path string, project_name string, module string) error {

		configDir := path + "/" + sandbox.ProjectName + "Config"
		if sandbox.Deps.IoLib.IsDir(configDir) {
			return sandbox.Deps.Errorf("the %s folder already exists", sandbox.ProjectName)
		}

		project_specs_dest := configDir + "/project.yaml"
		project_specs := sandbox.Deps.SerializeLib.CreateObject()

		project_specs.AddItemToObject("name", project_name)
		project_specs.AddItemToObject("module", module)
		project_specs.AddItemToObject("description", "Place the Project Description in "+project_specs_dest)
		res := sandbox.Deps.SerializeLib.SerializeToYaml(project_specs)
		sandbox.Deps.IoLib.WriteFile(project_specs_dest, []byte(res))

		sandbox.Deps.Printf("started with path %s \n", path)
		return nil
	}
}
