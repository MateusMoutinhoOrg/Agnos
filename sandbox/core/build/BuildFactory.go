package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// BuildFactory returns the closure that fills api.CoreApi.Build.
func BuildFactory(sandbox *api.SandBox) func(props api.BuildProps) error {
	return func(props api.BuildProps) error {

		sandbox.Deps.Printf("build started with path %s \n", props.Path)
		return nil
	}
}
