package cli_purge

import (
	"sort"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// cliDirs are the directories the cli layer owns whole. The asset group only
// names the files it installs, so removing those one by one would leave the
// generated neighbours behind — a command's entries.go with no entries.yaml
// and no handler.go next to it. The cli layer is generated from end to end,
// so purging it means dropping these directories entirely.
var cliDirs = []string{
	"sandbox/internal/cli",
	"sandbox/internal/commands",
}

// CliPurgeInternal removes from the target project every file that the "cli"
// asset group would have installed, at the path it holds inside that group,
// plus the directories the cli layer owns whole, then drops any directory the
// removal left empty. The deps the cli layer pulled in (sandbox/deps/argvdeps,
// sandbox/deps/std) are deliberately left in place: other code may use them.
func CliPurgeInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Log("cli-purge started with path %s \n", path)

	files, err := deps.Embeddeps.ListFilesRecursively("cli")
	if err != nil {
		return err
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	for _, dir := range cliDirs {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
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
