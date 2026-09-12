package list_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func ListDeps(sandbox *api.Sandbox, path string) ([]api.DepInfo, error) {
	io := smartio.New(sandbox, path, config.ProjectName)
	return ListDepsInternal(sandbox, io, path)
}
