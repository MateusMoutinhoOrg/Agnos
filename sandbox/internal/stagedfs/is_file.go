package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func IsFile(sandbox *api.Sandbox, io *StagedFS, path string) bool {
	p := processInputPath(io, path)
	if isPendingRemoval(sandbox, io, p) {
		return false
	}
	return sandbox.Deps.IoDeps.IsFile(rootedPath(sandbox, io, p))
}
