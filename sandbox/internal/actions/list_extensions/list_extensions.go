package list_extensions

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

func ListExtensions(sandbox *api.Sandbox, props api.ListExtensionsProps) ([]api.ExtensionInfo, error) {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return ListExtensionsInternal(sandbox, io, props.Path)
}
