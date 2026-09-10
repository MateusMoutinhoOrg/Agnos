package dep_remove

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// DepRemoveInternal uninstalls one dep: every adapter whose declaration names
// it, then the contract itself. An adapter is a directory, so removing one is
// removing adapters/libs/<adapter>/ — plus, for an adapter of the embedded
// catalog, whatever that catalog installs outside its own package.
func DepRemoveInternal(deps *deps.Deps, io *smartio.SmartIO, path string, dep string) error {
	deps.Std.Log("dep-remove started with path %s dep %s \n", path, dep)

	contract := utils.ContractsDir + "/" + dep
	if !io.IsDir(contract) {
		return deps.Std.Errorf("dep %q is not installed", dep)
	}

	for _, adapter := range utils.AdaptersFillingDep(deps, io, dep) {
		if err := RemoveAdapter(deps, io, adapter); err != nil {
			return err
		}
	}

	removeTree(deps, io, []string{contract})
	return nil
}

// RemoveAdapter drops one installed adapter: its enrollment in every
// available, its package directory, the extra files the embedded catalog
// installs alongside it, and the require its declaration pins.
func RemoveAdapter(deps *deps.Deps, io *smartio.SmartIO, adapter string) error {

	adapter_conf, err := utils.LoadAdapterConf(deps, io, adapter)
	if err != nil {
		return err
	}

	if err := utils.UnenrollAdapter(deps, io, adapter); err != nil {
		return err
	}

	removed := []string{utils.AdapterDir(adapter)}
	removed = append(removed, catalogExtras(deps, adapter)...)
	removeTree(deps, io, removed)

	module, _, ok := adapter_conf.ModuleSpec()
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

// catalogExtras returns the files assets/adapterlist/<adapter> installs
// outside the adapter's own package — the embed directive of `embeddeps` is
// one — so uninstalling takes back everything installing wrote. An adapter
// with no catalog entry (a generated shim) has none.
func catalogExtras(deps *deps.Deps, adapter string) []string {
	files, err := deps.Embeddeps.ListFilesRecursively(utils.AdapterlistGroup + "/" + adapter)
	if err != nil {
		return nil
	}

	var extras []string
	for _, file := range files {
		if file == utils.AdapterConfFile {
			continue
		}
		if deps.Stringsdeps.HasPrefix(file, utils.AdaptersDir+"/") {
			continue
		}
		extras = append(extras, file)
	}

	return extras
}

// removeTree removes each given path and then every directory the removal
// left empty, so the build collectors stop enumerating it.
func removeTree(deps *deps.Deps, io *smartio.SmartIO, paths []string) {
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
