package dep_install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/depsversionconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// DepInstallInternal renders every embedded asset under assets/deplist/<dep>
// into the target project at the path it holds inside that dep, using the
// same Module variable the build step derives from go.mod.
func DepInstallInternal(deps *deps.Deps, io *smartio.SmartIO, path string, dep string) error {
	deps.Std.Printf("dep-install started with path %s dep %s \n", path, dep)

	group := "deplist/" + dep

	files, err := deps.Embeddeps.ListFilesRecursively(group)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return deps.Std.Errorf("unknown dep %q", dep)
	}

	gomod, err := io.ReadFile(path + "/go.mod")
	if err != nil {
		return err
	}
	module_conf, err := moduleconf.New(deps, string(gomod))
	if err != nil {
		return err
	}

	if err := syncGoMod(deps, io, path, dep, module_conf); err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module": module_conf.Module,
	}

	return utils.RenderGroup(deps, io, group, vars)
}

// syncGoMod adds the require entry pinned for dep in assets/depsversion.yaml
// to the target project's go.mod, if that dep is listed there. Deps that
// bundle only sandbox-copy code (no external module) are absent from the file
// and leave go.mod untouched.
func syncGoMod(deps *deps.Deps, io *smartio.SmartIO, path string, dep string, module_conf *moduleconf.ModuleConf) error {
	versions_content, err := deps.Embeddeps.ReadFile("depsversion.yaml")
	if err != nil {
		return err
	}
	versions_conf, err := depsversionconf.New(deps, string(versions_content))
	if err != nil {
		return err
	}

	module, version, ok := versions_conf.Get(dep)
	if !ok {
		return nil
	}

	module_conf.AddRequire(module + " " + version)
	return io.WriteFileOverwrite(path+"/go.mod", []byte(module_conf.Render()))
}
