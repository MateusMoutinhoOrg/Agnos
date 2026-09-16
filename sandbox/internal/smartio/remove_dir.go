package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func RemoveDir(sandbox *api.Sandbox, io *SmartIO, path string) {
	p := processInputPath(io, path)
	io.PendingRemoveDirs = append(io.PendingRemoveDirs, p)
}
