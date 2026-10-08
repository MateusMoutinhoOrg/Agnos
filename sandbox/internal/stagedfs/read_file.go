package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ReadFile(sandbox *api.Sandbox, io *StagedFS, path string) ([]byte, error) {
	p := processInputPath(io, path)

	if isPendingRemoval(sandbox, io, p) {
		return nil, sandbox.Deps.StdDeps.Errorf("file %q does not exist", p)
	}

	if content, ok := io.Transactions[p]; ok {
		return content, nil
	}

	return sandbox.Deps.IoDeps.ReadFile(rootedPath(sandbox, io, p))
}
