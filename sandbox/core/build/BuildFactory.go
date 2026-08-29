package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

// BuildFactory returns the closure that fills core.CoreApi.Build.
func BuildFactory(sandbox *sandbox.SandBox) func(props core.BuildProps) error {
	return func(props core.BuildProps) error {

		sandbox.Deps.Printf("build started with path %s \n", props.Path)
		return nil
	}
}
