package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListFilesRecursively(sandbox *api.Sandbox, io *SmartIO, path string) []string {
	p := processInputPath(io, path)
	return filterPendingRemoved(sandbox, io, unrootedPaths(sandbox, io, sandbox.Deps.Iodeps.ListFilesRecursively(rootedPath(sandbox, io, p))))
}
