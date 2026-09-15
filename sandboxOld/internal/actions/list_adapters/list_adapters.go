package list_adapters

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func ListAdapters(sandbox *api.Sandbox, path string) ([]api.AdapterInfo, error) {
	io := smartio.New(sandbox, path, config.ProjectName)
	return ListAdaptersInternal(sandbox, io, path)
}
