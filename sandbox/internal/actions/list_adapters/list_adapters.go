package list_adapters

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func ListAdapters(deps *deps.Deps, path string) ([]api.AdapterInfo, error) {
	io := smartio.New(deps, path, config.ProjectName)
	return ListAdaptersInternal(deps, io, path)
}
