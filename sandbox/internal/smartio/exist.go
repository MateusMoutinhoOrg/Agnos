package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func Exist(deps *deps.Deps, io *SmartIO, path string) bool {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return false
	}
	return deps.IoLib.Exist(p)
}
