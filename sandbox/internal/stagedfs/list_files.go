package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListFiles(sandbox *api.Sandbox, io *StagedFS, path string) []string {
	p := processInputPath(io, path)
	return filterPendingRemoved(sandbox, io, unrootedPaths(sandbox, io, sandbox.Deps.IoDeps.ListFiles(rootedPath(sandbox, io, p))))
}
