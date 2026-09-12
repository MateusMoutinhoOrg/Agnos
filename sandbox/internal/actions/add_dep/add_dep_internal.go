package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddDepInternal installs one dep. An argument holding a "/" is the module
// path of another agnos repo and takes the remote route — the same
// disambiguation `go get` makes between a package name and a module path.
// Everything else is a name of the embedded catalog.
//
// The catalog route installs the contract under
// sandbox/deps/<dep>/ from assets/deplist/<dep>, and one adapter filling it
// from assets/adapterlist/<adapter> — the dep's declared default-adapter
// unless the caller names another. The two halves are separate catalogs, so
// the same contract can later be filled by a second implementation.
//
// The adapter is enrolled in every available: nothing else fills that field
// yet, and an available that leaves one empty is a nil func waiting to panic.
func AddDepInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.AddDepProps) error {
	if sandbox.Deps.Stringsdeps.Contains(props.Dep, "/") {
		return AddRemoteDepInternal(sandbox, io, props)
	}

	sandbox.Deps.Std.Log("add-dep started with path %s dep %s \n", props.Path, props.Dep)

	dep_conf, err := utils.LoadCatalogDepConf(sandbox, props.Dep)
	if err != nil {
		return err
	}

	adapter := props.Adapter
	if adapter == "" {
		adapter = dep_conf.DefaultAdapter
	}

	adapter_conf, err := utils.LoadCatalogAdapterConf(sandbox, adapter)
	if err != nil {
		return err
	}
	if adapter_conf.Dep != dep_conf.Name {
		return sandbox.Deps.Std.Errorf("adapter %q fills dep %q, not %q", adapter, adapter_conf.Dep, dep_conf.Name)
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module": module_conf.Module,
	}

	group := utils.DeplistGroup + "/" + dep_conf.Name
	if err := utils.RenderGroupExcept(sandbox, io, group, vars, []string{utils.DepConfFile}); err != nil {
		return err
	}

	if err := addAdapterAction.InstallAdapter(sandbox, io, adapter_conf, module_conf, vars); err != nil {
		return err
	}

	return utils.EnrollAdapter(sandbox, io, adapter_conf.Name)
}
