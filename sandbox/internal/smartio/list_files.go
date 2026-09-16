package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListFiles(sandbox *api.Sandbox, io *SmartIO, path string) []string {
	p := processInputPath(io, path)
	return filterPendingRemoved(sandbox, io, unrootedPaths(sandbox, io, sandbox.Deps.Iodeps.ListFiles(rootedPath(sandbox, io, p))))
}
