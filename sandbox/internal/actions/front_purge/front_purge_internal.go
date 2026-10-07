package front_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	backofficePurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/backoffice_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// FrontPurgeInternal removes from the target project every file that the
// "front" asset group would have installed, at the path it holds inside that
// group, plus the directories the front layer owns whole, then drops any
// directory the removal left empty.
//
// The frontend route is the one directory the front layer owns whole, looked
// up by name since it may have been moved to a folder: removing its files one
// by one would leave the generated new.go and entries.go behind with no
// route.yaml and no InternalPureHandler.go next to them, and its handler names
// the OpinatedAgnosFront lib, so leaving it behind would hand back a tree that
// does not serve what it says. What it serves does not go with it — assets/frontend/ is the
// project's own content, so front-init puts the route back over files that
// were never touched.
//
// The server layer is deliberately left in place, and so are the deps the
// front layer pulled in: other code may use them.
func FrontPurgeInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("front-purge started with path %s \n", path)

	// The backoffice stands on this layer, so it goes first, on this same open
	// SmartIO: a backoffice left on without it is a declaration verify
	// refuses. Its store on disk stays.
	has_backoffice, err := utils.ExtensionEnabled(sandbox, io, utils.ExtensionSandboxBackoffice)
	if err != nil {
		return err
	}
	if has_backoffice {
		if err := backofficePurgeAction.BackofficePurgeInternal(sandbox, io, path); err != nil {
			return err
		}
	}

	files, err := utils.ExtensionFiles(sandbox, utils.ExtensionSandboxFront)
	if err != nil {
		return err
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	for _, dir := range []string{utils.RouteDir(sandbox, io, utils.FrontendRouteName)} {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
	}

	utils.RemoveRetiredGenerated(sandbox, io, utils.ExtensionSandboxFront)

	sandbox.Deps.Std.Log("front-purge kept %s: every file there is yours \n", utils.FrontendDir)

	for _, dir := range ancestorDirs(sandbox, files) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	return utils.SetExtension(sandbox, io, utils.ExtensionSandboxFront, false)
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
	sandbox.Deps.Sortdeps.Slice(dirs, func(i int, j int) bool {
		return sandbox.Deps.Stringsdeps.Count(dirs[i], "/") > sandbox.Deps.Stringsdeps.Count(dirs[j], "/")
	})
	return dirs
}
