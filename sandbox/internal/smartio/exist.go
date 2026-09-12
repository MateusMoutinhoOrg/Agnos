package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func Exist(sandbox *api.Sandbox, io *SmartIO, path string) bool {
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
	return sandbox.Deps.Iodeps.Exist(rootedPath(sandbox, io, p))
}
