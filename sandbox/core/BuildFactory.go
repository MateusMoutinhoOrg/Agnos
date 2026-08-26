package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// BuildFactory returns the closure that fills api.CoreApi.Build.
func BuildFactory(sandbox *api.SandBox) func(path string) error {
	return func(path string) error {

		sandbox.Deps.Printf("build started with path %s \n", path)
		return nil
	}
}
