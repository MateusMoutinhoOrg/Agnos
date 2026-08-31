package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func IsDir(deps *deps.Deps, io *SmartIO, path string) bool {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return false
	}
	return deps.IoLib.IsDir(p)
}
