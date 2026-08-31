package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

// Build ...
func Build(sandbox *sandbox.SandBox, props core.BuildProps) error {
	sandbox.Deps.Printf("build started with path %s \n", props.Path)
	return nil
}
