package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func IsFile(sandbox *api.Sandbox, io *SmartIO, path string) bool {
	p, err := processInputPath(sandbox, io, path)
	if err != nil {
		return false
	}
	if isPendingRemoval(sandbox, io, p) {
		return false
	}
	return sandbox.Deps.Iodeps.IsFile(rootedPath(sandbox, io, p))
}
