package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveTree removes each given path and then every directory the removal left
// empty, so the build collectors stop enumerating it.
func RemoveTree(sandbox *api.Sandbox, io *smartio.SmartIO, paths []string) {
	for _, path := range paths {
		io.RemoveDir(path)
	}

	for _, dir := range ancestorDirs(sandbox, paths) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}
}

// ancestorDirs returns every directory that contains one of the given paths,
// deepest first, so an emptied child is removed before its parent is tested.
func ancestorDirs(sandbox *api.Sandbox, paths []string) []string {
	seen := map[string]bool{}
	var dirs []string
	for _, path := range paths {
		parts := sandbox.Deps.Stringsdeps.Split(path, "/")
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
