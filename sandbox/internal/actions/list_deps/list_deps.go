package list_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func ListDeps(deps *deps.Deps, path string) ([]api.DepInfo, error) {
	io := smartio.New(deps, path, config.ProjectName)
	return ListDepsInternal(deps, io, path)
}
