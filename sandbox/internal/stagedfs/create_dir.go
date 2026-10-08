package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func CreateDir(sandbox *api.Sandbox, io *StagedFS, path string) {
	p := processInputPath(io, path)
	io.PendingCreateDirs = append(io.PendingCreateDirs, p)
}
