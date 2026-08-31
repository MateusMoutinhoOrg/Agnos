package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func WriteFile(deps *deps.Deps, io *SmartIO, path string, content []byte) error {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return err
	}
	if deps.IoLib.Exist(p) {
		return deps.Errorf("file %q already exists", p)
	}
	io.Transactions[p] = content
	return nil
}
