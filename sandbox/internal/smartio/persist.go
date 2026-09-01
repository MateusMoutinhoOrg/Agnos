package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func Persist(deps *deps.Deps, io *SmartIO) error {
	for _, p := range io.PendingRemoveDirs {
		deps.IoLib.RemoveDir(p)
	}
	io.PendingRemoveDirs = nil

	for _, p := range io.PendingCreateDirs {
		deps.IoLib.CreateDir(p)
	}
	io.PendingCreateDirs = nil

	for p, content := range io.Transactions {
		err := deps.IoLib.WriteFile(p, content)
		if err != nil {
			return err
		}
	}
	io.Transactions = make(map[string][]byte)
	return nil
}
