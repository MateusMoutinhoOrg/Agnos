package server_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// serverDirs are the directories the server layer owns whole. The asset group
// only names the files it installs, so removing those one by one would leave
// the generated neighbours behind — a route's entries.go with no route.yaml
// and no handler.go next to it. The server layer is generated from end to end,
// so purging it means dropping these directories entirely.
//
// sandbox/internal/commands/start_server goes with them: it is the entry point
// server-init writes, and it would not compile without the layer it starts.
var serverDirs = []string{
	"sandbox/internal/server",
	"sandbox/internal/routes",
	"sandbox/internal/routeio",
	"sandbox/internal/commands/start_server",
}

// ServerPurgeInternal removes from the target project every file that the
// "server" asset group would have installed, at the path it holds inside that
// group, plus the directories the server layer owns whole, then drops any
// directory the removal left empty.
//
// The cli layer is deliberately left in place: server-init may have installed
// it, but a cli, once there, is the project's. The deps the server layer
// pulled in are left too — other code may use them.
func ServerPurgeInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Log("server-purge started with path %s \n", path)

	files, err := deps.Embeddeps.ListFilesRecursively("server")
	if err != nil {
		return err
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	for _, dir := range serverDirs {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
	}

	for _, dir := range ancestorDirs(deps, files) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	return nil
}

// ancestorDirs returns every directory that contains one of the given files,
// deepest first, so an emptied child is removed before its parent is tested.
func ancestorDirs(deps *deps.Deps, files []string) []string {
	seen := map[string]bool{}
	var dirs []string
	for _, file := range files {
		parts := deps.Stringsdeps.Split(file, "/")
		for i := 1; i < len(parts); i++ {
			dir := deps.Stringsdeps.Join(parts[:i], "/")
			if !seen[dir] {
				seen[dir] = true
				dirs = append(dirs, dir)
			}
		}
	}
	deps.Sortdeps.Slice(dirs, func(i, j int) bool {
		return deps.Stringsdeps.Count(dirs[i], "/") > deps.Stringsdeps.Count(dirs[j], "/")
	})
	return dirs
}
