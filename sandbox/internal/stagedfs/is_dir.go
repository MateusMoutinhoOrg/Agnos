package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func IsDir(sandbox *api.Sandbox, io *StagedFS, path string) bool {
	p := processInputPath(io, path)
	if isPendingRemoval(sandbox, io, p) {
		return false
	}
	if isPendingCreate(io, p) {
		return true
	}
	return sandbox.Deps.IoDeps.IsDir(rootedPath(sandbox, io, p))
}
