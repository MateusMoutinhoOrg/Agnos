package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func WriteFile(sandbox *api.Sandbox, io *StagedFS, path string, content []byte) error {
	p := processInputPath(io, path)
	io.Transactions[p] = content
	return nil
}
