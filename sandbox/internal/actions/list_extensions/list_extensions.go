package list_extensions

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func ListExtensions(sandbox *api.Sandbox, path string) ([]api.ExtensionInfo, error) {
	io := smartio.New(sandbox, path, config.ProjectName)
	return ListExtensionsInternal(sandbox, io, path)
}
