package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListFiles(sandbox *api.Sandbox, io *SmartIO, path string) []string {
	p, err := processInputPath(sandbox, io, path)
	if err != nil {
		return nil
	}
	return filterPendingRemoved(sandbox, io, filterIgnored(io, unrootedPaths(sandbox, io, sandbox.Deps.Iodeps.ListFiles(rootedPath(sandbox, io, p)))))
}
