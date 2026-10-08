package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListAllRecursively(sandbox *api.Sandbox, io *StagedFS, path string) []string {
	p := processInputPath(io, path)
	return filterPendingRemoved(sandbox, io, unrootedPaths(sandbox, io, sandbox.Deps.IoDeps.ListAllRecursively(rootedPath(sandbox, io, p))))
}
