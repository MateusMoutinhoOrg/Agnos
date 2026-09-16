package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func WriteFileOverwrite(sandbox *api.Sandbox, io *SmartIO, path string, content []byte) error {
	p := processInputPath(io, path)
	io.Transactions[p] = content
	return nil
}
