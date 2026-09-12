package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func RemoveDir(sandbox *api.Sandbox, io *SmartIO, path string) {
	p, err := processInputPath(sandbox, io, path)
	if err != nil {
		return
	}
	io.PendingRemoveDirs = append(io.PendingRemoveDirs, p)
}
