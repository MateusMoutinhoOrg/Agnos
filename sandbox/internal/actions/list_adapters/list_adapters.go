package list_adapters

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

func ListAdapters(sandbox *api.Sandbox, props api.ListAdaptersProps) ([]api.AdapterInfo, error) {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return ListAdaptersInternal(sandbox, io, props.Path)
}
