package dep_install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
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

	vars := map[string]interface{}{
		"Module": module_conf.Module,
	}

	return utils.RenderGroup(deps, io, group, vars)
}
