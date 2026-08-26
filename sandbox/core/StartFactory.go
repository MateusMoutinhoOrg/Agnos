package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func StartFactory(sandbox *api.SandBox) func(path string) error {
	return func(path string) error {

		sandbox.Deps.Printf("started with path %s \n", path)
		return nil
	}
}
