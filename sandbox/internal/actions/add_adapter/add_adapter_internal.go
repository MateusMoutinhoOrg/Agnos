package add_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddAdapterInternal installs one further implementation of a contract the
// project already has. The contract is the prerequisite, not the other way
// round: an adapter with no contract to fill is a package nothing can bind.
//
// Installing does not by itself change what any available binds — an available
// fills each field exactly once, so a second implementation only wins where it
// is asked to. props.Available names the one that switches to it.
func AddAdapterInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.AddAdapterProps) error {
	sandbox.Deps.Std.Log("add-adapter started with path %s adapter %s \n", props.Path, props.Adapter)

	adapter_conf, err := utils.LoadCatalogAdapterConf(sandbox, props.Adapter)
	if err != nil {
		return err
	}

	if !io.IsDir(utils.ContractsDir + "/" + adapter_conf.Dep) {
		return sandbox.Deps.Std.Errorf("adapter %q fills dep %q, which is not installed (run `agnos add-dep %s --adapter %s`)",
			adapter_conf.Name, adapter_conf.Dep, adapter_conf.Dep, adapter_conf.Name)
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	if err := InstallAdapter(sandbox, io, adapter_conf, module_conf, map[string]interface{}{
		"Module": module_conf.Module,
	}); err != nil {
		return err
	}

	if props.Available == "" {
		return nil
	}

	return utils.SelectAdapter(sandbox, io, props.Available, adapter_conf.Name)
}

// InstallAdapter renders one adapter of the embedded catalog into the target
// project and writes its declaration beside the package, so the installed tree
// answers by itself which dep that adapter fills. The versioned module the
// adapter imports — "" for one that needs nothing beyond the stdlib — is
// pinned in go.mod here, where the import actually is, rather than under the
// name of the contract.
//
// It enrolls the adapter nowhere: who binds it is the caller's decision, and
// the two callers make it differently — `add-dep` enrolls in every available
// because nothing else fills that field yet, `add-adapter` in the one named.
func InstallAdapter(sandbox *api.Sandbox, io *smartio.SmartIO, adapter_conf *adapterconf.AdapterConf, module_conf *moduleconf.ModuleConf, vars map[string]interface{}) error {

	group := utils.AdapterlistGroup + "/" + adapter_conf.Name
	if err := utils.RenderGroupExcept(sandbox, io, group, vars, []string{utils.AdapterConfFile}); err != nil {
		return err
	}

	if err := io.WriteFileOverwrite(utils.AdapterConfPath(adapter_conf.Name), []byte(adapter_conf.Render())); err != nil {
		return err
	}

	module, version, ok := adapter_conf.ModuleSpec()
	if !ok {
		return nil
	}

	module_conf.AddRequire(module + " " + version)
	return io.WriteFileOverwrite("go.mod", []byte(module_conf.Render()))
}
