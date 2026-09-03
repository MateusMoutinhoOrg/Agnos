package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"

func ReadFile(deps *deps.Deps, io *SmartIO, path string) ([]byte, error) {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return nil, err
	}

	if isPendingRemoval(io, p) {
		return nil, deps.Std.Errorf("file %q does not exist", p)
	}

	if content, ok := io.Transactions[p]; ok {
		return content, nil
	}

	return deps.Iodeps.ReadFile(rootedPath(io, p))
}
