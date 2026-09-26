package front_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// frontDirs are the directories the front layer owns whole: the file layer
// itself and the route serving the tree. The asset group only names the files
// it installs, so removing those one by one would leave the generated
// neighbours behind — the frontend route's new.go and entries.go with no
// route.yaml and no InternalPureHandler.go next to them.
var frontDirs = []string{
	"sandbox/internal/frontio",
	utils.RoutesDir + "/" + utils.FrontendRouteName,
}

// FrontPurgeInternal removes from the target project every file that the
// "front" asset group would have installed, at the path it holds inside that
// group, plus the directories the front layer owns whole, then drops any
// directory the removal left empty.
//
// The frontend route goes with the layer because its InternalPureHandler.go
// imports frontio, so leaving it behind would hand back a tree that does not
// compile. What it serves does not go with it — assets/frontend/ is the
// project's own content, so front-init puts the route back over files that
// were never touched.
//
// The server layer is deliberately left in place, and so are the deps the
// front layer pulled in: other code may use them.
func FrontPurgeInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("front-purge started with path %s \n", path)

	files, err := utils.ExtensionFiles(sandbox, utils.ExtensionSandboxFront)
	if err != nil {
		return err
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	for _, dir := range frontDirs {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
	}

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
