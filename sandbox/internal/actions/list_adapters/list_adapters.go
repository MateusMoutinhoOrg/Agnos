package list_adapters

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func ListAdapters(sandbox *api.Sandbox, path string) ([]api.AdapterInfo, error) {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	return ListAdaptersInternal(sandbox, io, path)
}
