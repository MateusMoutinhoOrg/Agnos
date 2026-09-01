package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func Build(deps *deps.Deps, path string) error {
	io := smartio.New(deps, path, config.ProjectName)
	err := BuildInternal(deps, io, path)
	if err != nil {
		return err
	}
	return io.Persist()
}
