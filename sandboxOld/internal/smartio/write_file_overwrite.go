package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func WriteFileOverwrite(sandbox *api.Sandbox, io *SmartIO, path string, content []byte) error {
	p, err := processInputPath(sandbox, io, path)
	if err != nil {
		return err
	}
	io.Transactions[p] = content
	return nil
}
