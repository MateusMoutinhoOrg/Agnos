package cli_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// cliDirs are the directories the cli layer owns whole. The asset group only
// names the files it installs, so removing those one by one would leave the
// generated neighbours behind — a command's new.go with no command.yaml
// and no handler.go next to it. The cli layer is generated from end to end,
// so purging it means dropping these directories entirely.
//
// sandbox/constructors/cli goes with them: it is what fills Sandbox.Cli, and
// it names the package being removed. It is written once and may since have
// been edited, so it is dropped with the layer it belongs to rather than left
// behind for sandbox/new.go to keep calling.
//
// docs/Commands is one of them for the same reason: the asset group installs
// its doc.md and doc.yaml, but the build writes one page per command beside
// them. Removing only the installed two would leave a directory of pages with
// no doc.yaml, which every later build reads as a doc that fails to load.
//
// sandbox/internal/cli is one of them too: its errors/ files are written once
// by the build and then the project's, but every one of them names api.Command,
// whose alias goes with the layer's api files, so leaving them behind hands
// back a tree that does not compile — and that no cli-init can bring back.
//
// sandbox/internal/commandprops is one of them: every command names it, so it
// goes with them.
var cliDirs = []string{
	utils.GeneratedDir + "/cli",
	"sandbox/internal/cli",
	"sandbox/internal/commands",
	utils.CommandPropsDir,
	"docs/Commands",
	utils.ConstructorDir("cli"),
}

// CliPurgeInternal removes from the target project every file that the "cli"
// asset group would have installed, at the path it holds inside that group,
// plus the directories the cli layer owns whole, then drops any directory the
// removal left empty. The deps the cli layer pulled in (sandbox/deps/argvdeps,
// sandbox/deps/stddeps) are deliberately left in place: other code may use them.
func CliPurgeInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("cli-purge started with path %s \n", path)

	files, err := utils.ExtensionFiles(sandbox, utils.ExtensionCli)
	if err != nil {
		return err
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	// sandbox/api/commandprops.go is where an older build declared
	// CommandProps; a tree no build has moved it out of yet still carries it.
	io.RemoveDir("sandbox/api/" + utils.CommandPropsFile)

	for _, dir := range cliDirs {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
	}

	utils.RemoveRetiredGenerated(sandbox, io, utils.ExtensionCli)

	for _, dir := range ancestorDirs(sandbox, files) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	return utils.SetExtension(sandbox, io, utils.ExtensionCli, false)
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
	sandbox.Deps.SortDeps.Slice(dirs, func(i, j int) bool {
		return sandbox.Deps.StringsDeps.Count(dirs[i], "/") > sandbox.Deps.StringsDeps.Count(dirs[j], "/")
	})
	return dirs
}
