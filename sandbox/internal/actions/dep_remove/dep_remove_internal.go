package dep_remove

import (
	"sort"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/depsversionconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// DepRemoveInternal removes from the target project every file that
// assets/deplist/<dep> would have installed, at the path it holds inside
// that dep, then drops any directory the removal left empty so the build
// collectors stop enumerating it.
func DepRemoveInternal(deps *deps.Deps, io *smartio.SmartIO, path string, dep string) error {
	deps.Std.Log("dep-remove started with path %s dep %s \n", path, dep)

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

	return syncGoMod(deps, io, path, dep)
}

// syncGoMod drops the require entry pinned for dep in assets/depsversion.yaml
// from the target project's go.mod, mirroring dep-install's addition. Deps
// absent from depsversion.yaml leave go.mod untouched.
func syncGoMod(deps *deps.Deps, io *smartio.SmartIO, path string, dep string) error {
	versions_content, err := deps.Embeddeps.ReadFile("depsversion.yaml")
	if err != nil {
		return err
	}
	versions_conf, err := depsversionconf.New(deps, string(versions_content))
	if err != nil {
		return err
	}

	module, _, ok := versions_conf.Get(dep)
	if !ok {
		return nil
	}

	module_conf, err := utils.LoadModuleConf(deps, io)
	if err != nil {
		return err
	}

	module_conf.RemoveRequire(module)
	return io.WriteFileOverwrite("go.mod", []byte(module_conf.Render()))
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
