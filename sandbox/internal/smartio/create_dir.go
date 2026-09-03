package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"

func CreateDir(deps *deps.Deps, io *SmartIO, path string) {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return
	}
	io.PendingCreateDirs = append(io.PendingCreateDirs, p)
}
