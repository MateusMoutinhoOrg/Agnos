package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListAllRecursively(sandbox *api.Sandbox, io *SmartIO, path string) []string {
	p := processInputPath(io, path)
	return filterPendingRemoved(sandbox, io, unrootedPaths(sandbox, io, sandbox.Deps.Iodeps.ListAllRecursively(rootedPath(sandbox, io, p))))
}
