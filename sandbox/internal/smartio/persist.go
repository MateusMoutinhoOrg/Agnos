package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"

func Persist(deps *deps.Deps, io *SmartIO) error {
	for _, p := range io.PendingRemoveDirs {
		deps.Iodeps.RemoveDir(rootedPath(deps, io, p))
	}
	io.PendingRemoveDirs = nil

	for _, p := range io.PendingCreateDirs {
		deps.Iodeps.CreateDir(rootedPath(deps, io, p))
	}
	io.PendingCreateDirs = nil

	for p, content := range io.Transactions {
		err := deps.Iodeps.WriteFile(rootedPath(deps, io, p), content)
		if err != nil {
			return err
		}
	}
	io.Transactions = make(map[string][]byte)
	return nil
}
