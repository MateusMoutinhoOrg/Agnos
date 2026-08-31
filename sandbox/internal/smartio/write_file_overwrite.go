package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func WriteFileOverwrite(deps *deps.Deps, io *SmartIO, path string, content []byte) error {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return err
	}
	io.Transactions[p] = content
	return nil
}
