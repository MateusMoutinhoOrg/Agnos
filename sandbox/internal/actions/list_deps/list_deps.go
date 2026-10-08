package list_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

func ListDeps(sandbox *api.Sandbox, props api.ListDepsProps) ([]api.DepInfo, error) {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return ListDepsInternal(sandbox, io, props.Path)
}
