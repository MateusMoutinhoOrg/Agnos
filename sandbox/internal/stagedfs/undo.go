package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

// Undo puts back every path the last Persist changed, in the reverse order it
// changed them: the files it wrote, the directories it created, then what it
// removed. It is a no-op when nothing was persisted.
func Undo(sandbox *api.Sandbox, io *StagedFS) error {
	journal := io.Journal
	if journal == nil {
		return nil
	}
	io.Journal = nil

	for index := len(journal.Order) - 1; index >= 0; index-- {
		rooted := journal.Order[index]
		if journal.Missing[rooted] {
			sandbox.Deps.IoDeps.RemoveDir(rooted)
			continue
		}
		if err := sandbox.Deps.IoDeps.WriteFile(rooted, journal.Files[rooted]); err != nil {
			return err
		}
	}

	for index := len(journal.CreatedDirs) - 1; index >= 0; index-- {
		sandbox.Deps.IoDeps.RemoveDir(journal.CreatedDirs[index])
	}

	for _, dir := range journal.RemovedDirs {
		sandbox.Deps.IoDeps.CreateDir(dir)
	}
	for rooted, content := range journal.RemovedFiles {
		sandbox.Deps.IoDeps.CreateDir(sandbox.Deps.IoDeps.Dir(rooted))
		if err := sandbox.Deps.IoDeps.WriteFile(rooted, content); err != nil {
			return err
		}
	}
	return nil
}
