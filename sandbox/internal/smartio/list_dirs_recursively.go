package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"

func ListDirsRecursively(deps *deps.Deps, io *SmartIO, path string) []string {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return nil
	}
	return filterPendingRemoved(deps, io, filterIgnored(io, unrootedPaths(deps, io, deps.Iodeps.ListDirsRecursively(rootedPath(deps, io, p)))))
}
