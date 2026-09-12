package front_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// routesDir holds one declared route per sub-directory; the pages this purge
// drops are the ones with an html template beside them.
const routesDir = "sandbox/internal/routes"

// frontDirs are the directories the front layer owns whole: the render layer
// itself and the route serving the static tree. The asset group only names the
// files it installs, so removing those one by one would leave the generated
// neighbours behind — the static route's entries.go with no route.yaml and no
// handler.go next to it.
var frontDirs = []string{
	"sandbox/internal/pageio",
	"sandbox/internal/routes/static",
}

// FrontPurgeInternal removes from the target project every file that the
// "front" asset group would have installed, at the path it holds inside that
// group, plus the directories the front layer owns whole and the route package
// of every declared page, then drops any directory the removal left empty.
//
// A page's route goes with the layer for the same reason server-purge drops
// sandbox/internal/routes whole: its handler.go imports pageio, so leaving it
// behind would hand back a tree that does not compile. What it renders does
// not go with it — assets/frontend/ is the project's own content, written by
// hand, so front-init followed by add-page puts the routes back over html that
// was never touched.
//
// The server layer is deliberately left in place, and so are the deps the
// front layer pulled in: other code may use them.
func FrontPurgeInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("front-purge started with path %s \n", path)

	files, err := sandbox.Deps.Embeddeps.ListFilesRecursively("front")
	if err != nil {
		return err
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	owned := make([]string, 0, len(frontDirs))
	owned = append(owned, frontDirs...)
	owned = append(owned, pageDirs(sandbox, io)...)

	for _, dir := range owned {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
	}

	sandbox.Deps.Std.Log("front-purge kept %s: pages, styles and scripts are yours \n", utils.FrontendDir)

	for _, dir := range ancestorDirs(sandbox, files) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	return nil
}

// pageDirs returns the route package of every declared page — every route with
// an html template beside it, the same test remove-page makes.
func pageDirs(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var dirs []string

	for _, dir := range io.ListDirs(routesDir) {
		name := lastSegmentOf(sandbox, dir)
		if name == "" || !utils.IsPage(sandbox, io, name) {
			continue
		}
		dirs = append(dirs, routesDir+"/"+name)
	}

	return dirs
}

// lastSegmentOf is the final slash-separated segment of a listed path, which
// is the route's directory name whether the listing came back rooted or not.
func lastSegmentOf(sandbox *api.Sandbox, path string) string {
	segments := sandbox.Deps.Stringsdeps.Split(path, "/")
	return segments[len(segments)-1]
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
