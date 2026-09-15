package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func Persist(sandbox *api.Sandbox, io *SmartIO) error {
	for _, p := range io.PendingRemoveDirs {
		sandbox.Deps.Iodeps.RemoveDir(rootedPath(sandbox, io, p))
	}
	io.PendingRemoveDirs = nil

	for _, p := range io.PendingCreateDirs {
		sandbox.Deps.Iodeps.CreateDir(rootedPath(sandbox, io, p))
	}
	io.PendingCreateDirs = nil

	for p, content := range io.Transactions {
		err := sandbox.Deps.Iodeps.WriteFile(rootedPath(sandbox, io, p), content)
		if err != nil {
			return err
		}
	}
	io.Transactions = make(map[string][]byte)
	return nil
}
