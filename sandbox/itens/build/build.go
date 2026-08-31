package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Build(deps deps.Deps, props api.BuildProps) error {
	deps.Printf("build started with path %s \n", props.Path)
	return nil
}
