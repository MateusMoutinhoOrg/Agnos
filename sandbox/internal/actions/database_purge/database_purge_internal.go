package database_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	backofficePurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/backoffice_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// databaseDirs are the directories the database layer owns whole. The asset
// group only names the files it installs, so removing those one by one would
// leave the generated neighbours behind — a database's api.go with no
// database.yaml next to it.
//
// docs/Databases is one of them: the doc group installs its doc.md and
// doc.yaml, but the build writes one page per database beside them. Removing
// only the installed two would leave a directory of pages with no doc.yaml,
// which every later build reads as a doc that fails to load.
//
// methods_custom.go goes with them. It is hand-written, and that is exactly
// why: it is Go in a package whose other three files are about to be gone, so
// leaving it behind would leave the project not compiling.
//
// So does the --database middleware: database-cli installs its command.yaml
// and handler.go, and the build writes new.go and input.go beside them.
var databaseDirs = []string{
	utils.DatabasesDir,
	utils.DocsDir + "/Databases",
	utils.DatabaseDirMiddleware,
}

// DatabasePurgeInternal removes from the target project every file that the
// database asset groups would have installed, at the path it holds inside each
// group, plus the directories the database layer owns whole, then drops any
// directory the removal left empty.
//
// The database dep is left in place: database-init installed it, but a
// contract, once there, is the project's — `remove-dep database` is what takes
// it away.
func DatabasePurgeInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("database-purge started with path %s \n", path)

	// The backoffice stands on this layer, so it goes first, on this same open
	// StagedFS: a backoffice left on without it is a declaration verify
	// refuses. Its store on disk stays.
	has_backoffice, err := utils.ExtensionEnabled(sandbox, io, utils.ExtensionBackoffice)
	if err != nil {
		return err
	}
	if has_backoffice {
		if err := backofficePurgeAction.BackofficePurgeInternal(sandbox, io, path); err != nil {
			return err
		}
	}

	files, err := utils.ExtensionFiles(sandbox, utils.ExtensionDatabase)
	if err != nil {
		return err
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	for _, dir := range databaseDirs {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
	}

	utils.RemoveRetiredGenerated(sandbox, io, utils.ExtensionDatabase)

	for _, dir := range ancestorDirs(sandbox, files) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	return utils.SetExtension(sandbox, io, utils.ExtensionDatabase, false)
}

// ancestorDirs returns every directory that contains one of the given files,
// deepest first, so an emptied child is removed before its parent is tested.
func ancestorDirs(sandbox *api.Sandbox, files []string) []string {
	seen := map[string]bool{}
	var dirs []string
	for _, file := range files {
		parts := sandbox.Deps.StringsDeps.Split(file, "/")
		for i := 1; i < len(parts); i++ {
			dir := sandbox.Deps.StringsDeps.Join(parts[:i], "/")
			if !seen[dir] {
				seen[dir] = true
				dirs = append(dirs, dir)
			}
		}
	}
	sandbox.Deps.SortDeps.Slice(dirs, func(i, j int) bool {
		return sandbox.Deps.StringsDeps.Count(dirs[i], "/") > sandbox.Deps.StringsDeps.Count(dirs[j], "/")
	})
	return dirs
}
