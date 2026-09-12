package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func CreateDir(sandbox *api.Sandbox, io *SmartIO, path string) {
	p, err := processInputPath(sandbox, io, path)
	if err != nil {
		return
	}
	io.PendingCreateDirs = append(io.PendingCreateDirs, p)
}
