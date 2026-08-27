package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func StartFactory(sandbox *api.SandBox) func(props api.StartProps) error {

	return func(props api.StartProps) error {

		//configDir := props.Path + "/" + sandbox.ProjectName + "Config"

		sandbox.Deps.Printf("started with path %s \n", props.Path)
		return nil
	}
}
