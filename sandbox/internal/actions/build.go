package actions

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Build(deps *deps.Deps, path string) error {
	deps.Printf("build started with path %s \n", path)
	return nil
}
