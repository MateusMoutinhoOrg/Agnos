package remove_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveDatabaseInternal deletes every file under
// sandbox/internal/databases/<name>/ plus the directory itself. A package
// carrying a methods_custom.go is refused: that file is the project's own, and
// nothing in the declaration says what it holds, so taking it away silently is
// the one removal agnos will not do on its own.
func RemoveDatabaseInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) error {
	if err := utils.ValidateDatabaseName(sandbox, name); err != nil {
		return err
	}

	dir := utils.DatabaseDir(sandbox, name)
	if !io.IsDir(dir) {
		return sandbox.Deps.StdDeps.Errorf("database %q not found", utils.DatabaseName(sandbox, name))
	}

	custom := utils.DatabaseCustomPath(sandbox, name)
	if io.IsFile(custom) {
		return sandbox.Deps.StdDeps.Errorf(
			"database %q carries %s, which is hand-written: delete that file first if it is really going",
			utils.DatabaseName(sandbox, name), custom)
	}

	sandbox.Deps.StdDeps.Logf("remove-database removing %s \n", dir)

	for _, file := range io.ListAllRecursively(dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(dir)
	return nil
}
