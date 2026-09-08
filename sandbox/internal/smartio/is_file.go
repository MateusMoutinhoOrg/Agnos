package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"

func IsFile(deps *deps.Deps, io *SmartIO, path string) bool {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return false
	}
	if isPendingRemoval(deps, io, p) {
		return false
	}
	return deps.Iodeps.IsFile(rootedPath(deps, io, p))
}
