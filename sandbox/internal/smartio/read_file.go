package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func ReadFile(deps *deps.Deps, io *SmartIO, path string) ([]byte, error) {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return nil, err
	}
	return deps.IoLib.ReadFile(p)
}
