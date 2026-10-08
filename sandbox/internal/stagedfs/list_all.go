package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListAll(sandbox *api.Sandbox, io *StagedFS, path string) []string {
	p := processInputPath(io, path)
	return filterPendingRemoved(sandbox, io, unrootedPaths(sandbox, io, sandbox.Deps.IoDeps.ListAll(rootedPath(sandbox, io, p))))
}
