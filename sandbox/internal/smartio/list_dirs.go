package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func ListDirs(deps *deps.Deps, io *SmartIO, path string) []string {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return nil
	}
	return filterPendingRemoved(io, filterIgnored(io, unrootedPaths(io, deps.Iodeps.ListDirs(rootedPath(io, p)))))
}
