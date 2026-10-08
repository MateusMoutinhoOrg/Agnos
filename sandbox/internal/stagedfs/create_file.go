package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func CreateFile(sandbox *api.Sandbox, io *StagedFS, path string, content []byte) error {
	p := processInputPath(io, path)
	if sandbox.Deps.IoDeps.Exists(rootedPath(sandbox, io, p)) {
		return sandbox.Deps.StdDeps.Errorf("file %q already exists", p)
	}
	io.Transactions[p] = content
	return nil
}
