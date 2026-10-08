package stagedfs

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

// Persist writes the pending transaction to disk: removals first, then
// creations, then every file. What each path held before is recorded in
// io.Journal first, so Undo can restore it.
func Persist(sandbox *api.Sandbox, io *StagedFS) error {
	journal := &Journal{Files: map[string][]byte{}, Missing: map[string]bool{}, RemovedFiles: map[string][]byte{}}
	io.Journal = journal

	for _, p := range io.PendingRemoveDirs {
		rooted := rootedPath(sandbox, io, p)
		recordRemoval(sandbox, journal, rooted)
		sandbox.Deps.IoDeps.RemoveDir(rooted)
	}
	io.PendingRemoveDirs = nil

	for _, p := range io.PendingCreateDirs {
		rooted := rootedPath(sandbox, io, p)
		if !sandbox.Deps.IoDeps.Exists(rooted) {
			journal.CreatedDirs = append(journal.CreatedDirs, rooted)
		}
		sandbox.Deps.IoDeps.CreateDir(rooted)
	}
	io.PendingCreateDirs = nil

	for p, content := range io.Transactions {
		rooted := rootedPath(sandbox, io, p)
		recordFile(sandbox, journal, rooted)
		err := sandbox.Deps.IoDeps.WriteFile(rooted, content)
		if err != nil {
			return err
		}
	}
	io.Transactions = make(map[string][]byte)
	return nil
}

// recordFile keeps what a file held before its first write, or that it was
// not there.
func recordFile(sandbox *api.Sandbox, journal *Journal, rooted string) {
	if _, recorded := journal.Files[rooted]; recorded || journal.Missing[rooted] {
		return
	}
	journal.Order = append(journal.Order, rooted)
	if !sandbox.Deps.IoDeps.IsFile(rooted) {
		journal.Missing[rooted] = true
		recordMissingParent(sandbox, journal, rooted)
		return
	}
	content, err := sandbox.Deps.IoDeps.ReadFile(rooted)
	if err != nil {
		journal.Missing[rooted] = true
		return
	}
	journal.Files[rooted] = content
}

// recordRemoval keeps the whole contents of a directory — or a file — about to
// be removed.
func recordRemoval(sandbox *api.Sandbox, journal *Journal, rooted string) {
	if sandbox.Deps.IoDeps.IsFile(rooted) {
		if content, err := sandbox.Deps.IoDeps.ReadFile(rooted); err == nil {
			journal.RemovedFiles[rooted] = content
		}
		return
	}
	if !sandbox.Deps.IoDeps.IsDir(rooted) {
		return
	}
	journal.RemovedDirs = append(journal.RemovedDirs, rooted)
	journal.RemovedDirs = append(journal.RemovedDirs, sandbox.Deps.IoDeps.ListDirsRecursively(rooted)...)
	for _, file := range sandbox.Deps.IoDeps.ListFilesRecursively(rooted) {
		if content, err := sandbox.Deps.IoDeps.ReadFile(file); err == nil {
			journal.RemovedFiles[file] = content
		}
	}
}

// recordMissingParent keeps the outermost directory a write to rooted is about
// to create, since writing a file creates every directory above it.
func recordMissingParent(sandbox *api.Sandbox, journal *Journal, rooted string) {
	missing := ""
	for dir := sandbox.Deps.IoDeps.Dir(rooted); dir != "" && dir != "." && dir != "/"; dir = sandbox.Deps.IoDeps.Dir(dir) {
		if sandbox.Deps.IoDeps.Exists(dir) {
			break
		}
		missing = dir
	}
	if missing == "" {
		return
	}
	for _, created := range journal.CreatedDirs {
		if created == missing {
			return
		}
	}
	journal.CreatedDirs = append(journal.CreatedDirs, missing)
}
