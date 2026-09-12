package remove_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveAdapterInternal uninstalls one adapter and nothing else: the contract
// it filled stays, because a contract may have more than one implementation.
//
// It refuses two cases. An adapter some available still binds would leave that
// available with an unfilled field — a nil func waiting to panic — so the
// caller is told to point that available at another adapter first. An adapter
// the generator wrote is half of a dep copied from a remote repo; the contract
// beside it is the other half, and `remove-dep` is what owns the pair.
func RemoveAdapterInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string, adapter string) error {
	sandbox.Deps.Std.Log("remove-adapter started with path %s adapter %s \n", path, adapter)

	adapter_conf, err := utils.LoadAdapterConf(sandbox, io, adapter)
	if err != nil {
		return sandbox.Deps.Std.Errorf("adapter %q is not installed", adapter)
	}

	if adapter_conf.Origin == adapterconf.OriginGenerated {
		return sandbox.Deps.Std.Errorf("adapter %q was generated as the shim of dep %q: remove the pair with `agnos remove-dep %s`",
			adapter, adapter_conf.Dep, adapter_conf.Dep)
	}

	if binding := utils.AvailablesBinding(sandbox, io, adapter); len(binding) > 0 {
		return sandbox.Deps.Std.Errorf("adapter %q is the only one filling deps.%s in %s: point it at another adapter first (`agnos set-adapter %s <adapter>`)",
			adapter, utils.DepField(sandbox, adapter_conf.Dep),
			sandbox.Deps.Stringsdeps.Join(binding, ", "), adapter_conf.Dep)
	}

	return Uninstall(sandbox, io, adapter_conf)
}

// Uninstall drops one adapter's files and its require, with no question asked
// about who binds it. `remove-dep` reaches it directly: it is taking the whole
// pair, so the refusals of RemoveAdapterInternal do not apply.
func Uninstall(sandbox *api.Sandbox, io *smartio.SmartIO, adapter_conf *adapterconf.AdapterConf) error {

	if err := utils.UnenrollAdapter(sandbox, io, adapter_conf.Name); err != nil {
		return err
	}

	removed := []string{utils.AdapterDir(adapter_conf.Name)}
	removed = append(removed, catalogExtras(sandbox, adapter_conf.Name)...)
	utils.RemoveTree(sandbox, io, removed)

	module, _, ok := adapter_conf.ModuleSpec()
	if !ok {
		return nil
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
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
func catalogExtras(sandbox *api.Sandbox, adapter string) []string {
	files, err := sandbox.Deps.Embeddeps.ListFilesRecursively(utils.AdapterlistGroup + "/" + adapter)
	if err != nil {
		return nil
	}

	var extras []string
	for _, file := range files {
		if file == utils.AdapterConfFile {
			continue
		}
		if sandbox.Deps.Stringsdeps.HasPrefix(file, utils.AdaptersDir+"/") {
			continue
		}
		extras = append(extras, file)
	}

	return extras
}
