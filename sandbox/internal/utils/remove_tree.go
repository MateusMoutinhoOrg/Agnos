package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveTree removes each given path and then every directory the removal left
// empty, so the build collectors stop enumerating it.
func RemoveTree(deps *deps.Deps, io *smartio.SmartIO, paths []string) {
	for _, path := range paths {
		io.RemoveDir(path)
	}

	for _, dir := range ancestorDirs(deps, paths) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}
}

// ancestorDirs returns every directory that contains one of the given paths,
// deepest first, so an emptied child is removed before its parent is tested.
func ancestorDirs(deps *deps.Deps, paths []string) []string {
	seen := map[string]bool{}
	var dirs []string
	for _, path := range paths {
		parts := deps.Stringsdeps.Split(path, "/")
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
