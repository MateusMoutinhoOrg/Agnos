package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func RemoveDir(deps *deps.Deps, io *SmartIO, path string) {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return
	}
	io.PendingRemoveDirs = append(io.PendingRemoveDirs, p)
}
