package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func WriteFile(sandbox *api.Sandbox, io *SmartIO, path string, content []byte) error {
	p := processInputPath(io, path)
	if sandbox.Deps.Iodeps.Exist(rootedPath(sandbox, io, p)) {
		return sandbox.Deps.Std.Errorf("file %q already exists", p)
	}
	io.Transactions[p] = content
	return nil
}
