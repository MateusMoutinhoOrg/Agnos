package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func StartFactory(sandbox *api.SandBox) func(props api.StartProps) error {

	return func(props api.StartProps) error {

		configDir := props.Path + "/" + sandbox.ProjectName + "Config"
		if !props.Force {
			if sandbox.Deps.IoLib.IsDir(configDir) {
				return sandbox.Deps.Errorf("the %s folder already exists", sandbox.ProjectName)
			}
		}
		err := createProjectSpecs(sandbox, props, configDir)
		if err != nil {
			return err
		}

		sandbox.Deps.Printf("started with path %s \n", props.Path)
		return nil
	}
}
