package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func Exists(sandbox *api.Sandbox, io *StagedFS, path string) bool {
	p := processInputPath(io, path)
	if isPendingRemoval(sandbox, io, p) {
		return false
	}
	if isPendingCreate(io, p) {
		return true
	}
	return sandbox.Deps.IoDeps.Exists(rootedPath(sandbox, io, p))
}
