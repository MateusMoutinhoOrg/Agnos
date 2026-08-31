package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

type BuildProps struct {
	Path    string
	Project string
	Force   bool
}

func Build(deps deps.Deps, props BuildProps) error {
	deps.Printf("build started with path %s \n", props.Path)
	return nil
}
