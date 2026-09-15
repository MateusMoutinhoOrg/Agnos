package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ListAll(sandbox *api.Sandbox, io *SmartIO, path string) []string {
	p, err := processInputPath(sandbox, io, path)
	if err != nil {
		return nil
	}
	return filterPendingRemoved(sandbox, io, filterIgnored(io, unrootedPaths(sandbox, io, sandbox.Deps.Iodeps.ListAll(rootedPath(sandbox, io, p)))))
}
