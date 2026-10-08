package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func RemoveDir(sandbox *api.Sandbox, io *StagedFS, path string) {
	p := processInputPath(io, path)
	io.PendingRemoveDirs = append(io.PendingRemoveDirs, p)
}
