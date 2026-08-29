package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

// BuildFactory returns the closure that fills lib.CoreApi.Build.
func BuildFactory(sandbox *lib.SandBox) func(props lib.BuildProps) error {
	return func(props lib.BuildProps) error {

		sandbox.Deps.Printf("build started with path %s \n", props.Path)
		return nil
	}
}
