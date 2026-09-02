package cli_purge

import (
	"sort"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// CliPurgeInternal removes from the target project every file that the "cli"
// asset group would have installed, at the path it holds inside that group,
// then drops any directory the removal left empty.
func CliPurgeInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Printf("cli-purge started with path %s \n", path)

	files, err := deps.Embeddeps.ListFilesRecursively("cli")
	if err != nil {
		return err
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
