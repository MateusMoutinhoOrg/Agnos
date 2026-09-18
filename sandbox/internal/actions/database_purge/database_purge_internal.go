package database_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// databaseDirs are the directories the database layer owns whole. The asset
// group only names the files it installs, so removing those one by one would
// leave the generated neighbours behind — a database's api.go with no
// specs.yaml next to it.
//
// docs/Databases is one of them: the doc group installs its doc.md and
// props.yaml, but the build writes one page per database beside them. Removing
// only the installed two would leave a directory of pages with no props.yaml,
// which every later build reads as a doc that fails to load.
//
// methods_custom.go goes with them. It is hand-written, and that is exactly
// why: it is Go in a package whose other three files are about to be gone, so
// leaving it behind would leave the project not compiling.
var databaseDirs = []string{
	utils.DatabasesDir,
	utils.DocsDir + "/Databases",
}

// DatabasePurgeInternal removes from the target project every file that the
// database asset groups would have installed, at the path it holds inside each
// group, plus the directories the database layer owns whole, then drops any
// directory the removal left empty.
//
// The database dep is left in place: database-init installed it, but a
// contract, once there, is the project's — `remove-dep database` is what takes
// it away.
func DatabasePurgeInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("database-purge started with path %s \n", path)

	files, err := utils.ExtensionFiles(sandbox, utils.ExtensionSandboxDatabase)
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

	for _, dir := range ancestorDirs(sandbox, files) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	return utils.SetExtension(sandbox, io, utils.ExtensionSandboxDatabase, false)
}

// ancestorDirs returns every directory that contains one of the given files,
// deepest first, so an emptied child is removed before its parent is tested.
func ancestorDirs(sandbox *api.Sandbox, files []string) []string {
	seen := map[string]bool{}
	var dirs []string
	for _, file := range files {
		parts := sandbox.Deps.Stringsdeps.Split(file, "/")
		for i := 1; i < len(parts); i++ {
			dir := sandbox.Deps.Stringsdeps.Join(parts[:i], "/")
			if !seen[dir] {
				seen[dir] = true
				dirs = append(dirs, dir)
			}
		}
	}
	sandbox.Deps.Sortdeps.Slice(dirs, func(i, j int) bool {
		return sandbox.Deps.Stringsdeps.Count(dirs[i], "/") > sandbox.Deps.Stringsdeps.Count(dirs[j], "/")
	})
	return dirs
}
