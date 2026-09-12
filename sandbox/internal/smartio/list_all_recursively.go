package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListAllRecursively(sandbox *api.Sandbox, io *SmartIO, path string) []string {
	p, err := processInputPath(sandbox, io, path)
	if err != nil {
		return nil
	}
	return filterPendingRemoved(sandbox, io, filterIgnored(io, unrootedPaths(sandbox, io, sandbox.Deps.Iodeps.ListAllRecursively(rootedPath(sandbox, io, p)))))
}
