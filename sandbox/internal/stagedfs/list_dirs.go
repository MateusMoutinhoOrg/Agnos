package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListDirs(sandbox *api.Sandbox, io *SmartIO, path string) []string {
	p := processInputPath(io, path)
	return filterPendingRemoved(sandbox, io, unrootedPaths(sandbox, io, sandbox.Deps.Iodeps.ListDirs(rootedPath(sandbox, io, p))))
}
