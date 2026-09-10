package dep_install

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// DepInstallInternal installs one dep of the embedded catalog: the contract
// under sandbox/deps/<dep>/ from assets/deplist/<dep>, and one adapter filling
// it from assets/adapterlist/<adapter> — the dep's declared default-adapter
// unless the caller names another. The two halves are separate catalogs, so
// the same contract can later be filled by a second implementation.
func DepInstallInternal(deps *deps.Deps, io *smartio.SmartIO, props api.DepInstallProps) error {
	deps.Std.Log("dep-install started with path %s dep %s \n", props.Path, props.Dep)

	dep_conf, err := utils.LoadCatalogDepConf(deps, props.Dep)
	if err != nil {
		return err
	}

	adapter := props.Adapter
	if adapter == "" {
		adapter = dep_conf.DefaultAdapter
	}

	adapter_conf, err := utils.LoadCatalogAdapterConf(deps, adapter)
	if err != nil {
		return err
	}
	if adapter_conf.Dep != dep_conf.Name {
		return deps.Std.Errorf("adapter %q fills dep %q, not %q", adapter, adapter_conf.Dep, dep_conf.Name)
	}

	module_conf, err := utils.LoadModuleConf(deps, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module": module_conf.Module,
	}

	group := utils.DeplistGroup + "/" + dep_conf.Name
	if err := utils.RenderGroupExcept(deps, io, group, vars, []string{utils.DepConfFile}); err != nil {
		return err
	}

	return InstallAdapter(deps, io, adapter_conf, module_conf, vars)
}
