package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func CreateDir(deps *deps.Deps, io *SmartIO, path string) {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return
	}
	deps.IoLib.CreateDir(p)
}
