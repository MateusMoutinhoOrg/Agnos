package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func ReadFile(sandbox *api.Sandbox, io *SmartIO, path string) ([]byte, error) {
	p, err := processInputPath(sandbox, io, path)
	if err != nil {
		return nil, err
	}

	if isPendingRemoval(sandbox, io, p) {
		return nil, sandbox.Deps.Std.Errorf("file %q does not exist", p)
	}

	if content, ok := io.Transactions[p]; ok {
		return content, nil
	}

	return sandbox.Deps.Iodeps.ReadFile(rootedPath(sandbox, io, p))
}
