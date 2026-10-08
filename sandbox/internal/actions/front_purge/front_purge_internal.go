package front_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	backofficePurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/backoffice_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// FrontPurgeInternal removes from the target project every file that the
// "front" asset group would have installed, at the path it holds inside that
// group, plus the directories the front layer owns whole, then drops any
// directory the removal left empty.
//
// The front route is the one directory the front layer owns whole, looked
// up by name since it may have been moved to a folder: removing its files one
// by one would leave the generated new.go and input.go behind with no
// route.yaml and no handler.go next to them, and its handler names
// the OpinionatedAgnosFront lib, so leaving it behind would hand back a tree that
// does not serve what it says. What it serves does not go with it — assets/front/ is the
// project's own content, so front-init puts the route back over files that
// were never touched.
//
// The server layer is deliberately left in place, and so are the deps the
// front layer pulled in: other code may use them.
func FrontPurgeInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("front-purge started with path %s \n", path)

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

	files, err := utils.ExtensionFiles(sandbox, utils.ExtensionFront)
	if err != nil {
		return err
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	for _, dir := range []string{utils.RouteDir(sandbox, io, utils.FrontRouteName)} {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
	}

	utils.RemoveRetiredGenerated(sandbox, io, utils.ExtensionFront)

	sandbox.Deps.StdDeps.Logf("front-purge kept %s: every file there is yours \n", utils.FrontDir)

	for _, dir := range ancestorDirs(sandbox, files) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	return utils.SetExtension(sandbox, io, utils.ExtensionFront, false)
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
	sandbox.Deps.SortDeps.Slice(dirs, func(i int, j int) bool {
		return sandbox.Deps.StringsDeps.Count(dirs[i], "/") > sandbox.Deps.StringsDeps.Count(dirs[j], "/")
	})
	return dirs
}
