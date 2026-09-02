package dep_remove

import (
	"sort"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// DepRemoveInternal removes from the target project every file that
// assets/deplist/<dep> would have installed, at the path it holds inside
// that dep, then drops any directory the removal left empty so the build
// collectors stop enumerating it.
func DepRemoveInternal(deps *deps.Deps, io *smartio.SmartIO, path string, dep string) error {
	deps.Std.Printf("dep-remove started with path %s dep %s \n", path, dep)

	group := "deplist/" + dep

	files, err := deps.Embeddeps.ListFilesRecursively(group)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return deps.Std.Errorf("unknown dep %q", dep)
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	for _, dir := range ancestorDirs(files) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	return nil
}

// ancestorDirs returns every directory that contains one of the given files,
// deepest first, so an emptied child is removed before its parent is tested.
func ancestorDirs(files []string) []string {
	seen := map[string]bool{}
	var dirs []string
	for _, file := range files {
		parts := strings.Split(file, "/")
		for i := 1; i < len(parts); i++ {
			dir := strings.Join(parts[:i], "/")
			if !seen[dir] {
				seen[dir] = true
				dirs = append(dirs, dir)
			}
		}
	}
	sort.Slice(dirs, func(i, j int) bool {
		return strings.Count(dirs[i], "/") > strings.Count(dirs[j], "/")
	})
	return dirs
}
