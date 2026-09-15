package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func IsDir(sandbox *api.Sandbox, io *SmartIO, path string) bool {
	p, err := processInputPath(sandbox, io, path)
	if err != nil {
		return false
	}
	if isPendingRemoval(sandbox, io, p) {
		return false
	}
	if isPendingCreate(io, p) {
		return true
	}
	return sandbox.Deps.Iodeps.IsDir(rootedPath(sandbox, io, p))
}
