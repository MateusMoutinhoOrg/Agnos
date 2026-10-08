package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func CreateDir(sandbox *api.Sandbox, io *SmartIO, path string) {
	p := processInputPath(io, path)
	io.PendingCreateDirs = append(io.PendingCreateDirs, p)
}
