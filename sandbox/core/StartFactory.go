package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func StartFactory(sandbox *api.SandBox) func(path string) error {

	return func(path string) error {
		if sandbox.Deps.IoLib.IsDir(sandbox.ProjectName + "Config") {
			return sandbox.Deps.Errorf("the %s folder already exists", sandbox.ProjectName)
		}

		sandbox.Deps.Printf("started with path %s \n", path)
		return nil
	}
}
