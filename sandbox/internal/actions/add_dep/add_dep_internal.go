package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddDepInternal installs one dep. An argument holding a "/" is the module
// path of another agnos repo and takes the remote route — the same
// disambiguation `go get` makes between a package name and a module path.
// Everything else is a name of the embedded catalog.
//
// The catalog route installs the contract under
// sandbox/deps/<dep>/ from assets/dep-catalog/<dep>, and one adapter filling it
// from assets/adapter-catalog/<adapter> — the dep's declared default-adapter
// unless the caller names another. The two halves are separate catalogs, so
// the same contract can later be filled by a second implementation.
//
// The adapter is enrolled in every binding: nothing else fills that field
// yet, and a binding that leaves one empty is a nil func waiting to panic.
func AddDepInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.AddDepProps) error {
	// Installing a dep is asking for the dependency layer, so the mechanic
	// that renders sandbox/deps/deps.go is turned on here rather than being
	// inferred later from the directory this install is about to create.
	if err := utils.SetExtension(sandbox, io, utils.ExtensionDeps, true); err != nil {
		return err
	}

	if sandbox.Deps.StringsDeps.Contains(props.Dep, "/") {
		return AddRemoteDepInternal(sandbox, io, props)
	}

	sandbox.Deps.StdDeps.Logf("add-dep started with path %s dep %s \n", props.Path, props.Dep)

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
		return sandbox.Deps.StdDeps.Errorf("adapter %q fills dep %q, not %q", adapter, adapter_conf.Dep, dep_conf.Name)
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module": module_conf.Module,
	}

	group := utils.DepCatalogGroup + "/" + dep_conf.Name
	if err := utils.RenderGroupExcept(sandbox, io, group, vars, []string{utils.DepConfFile}); err != nil {
		return err
	}

	if err := addAdapterAction.InstallAdapter(sandbox, io, adapter_conf, module_conf, vars); err != nil {
		return err
	}

	return utils.EnrollAdapter(sandbox, io, adapter_conf.Name)
}
