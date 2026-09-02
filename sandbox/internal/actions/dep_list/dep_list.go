package dep_list

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func DepList(deps *deps.Deps, path string) ([]string, error) {
	io := smartio.New(deps, path, config.ProjectName)
	return DepListInternal(deps, io, path)
}
