package set_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetDepInternal moves one remote dep to another version of the module it was
// copied from: the contract is copied again and the shim regenerated, so a
// contract that changed shape is carried across in one step instead of being
// hand-patched. The module it came from is read back out of the shim's own
// declaration, which is where the install filed it.
//
// It is only for a dep whose adapter the generator wrote. A dep of the
// embedded catalog has no version of its own — its adapter's does, and that
// one moves with agnos itself.
func SetDepInternal(deps *deps.Deps, io *smartio.SmartIO, props api.SetDepProps) error {
	deps.Std.Log("set-dep started with path %s dep %s version %s \n", props.Path, props.Dep, props.Version)

	adapter_conf, err := utils.LoadAdapterConf(deps, io, props.Dep)
	if err != nil {
		return deps.Std.Errorf("dep %q has no generated adapter: only a dep copied from another repo has a version to set", props.Dep)
	}

	if adapter_conf.Origin != adapterconf.OriginGenerated {
		return deps.Std.Errorf("dep %q came from the embedded catalog, which carries no version of its own: change the adapter instead (`agnos set-adapter %s <adapter>`)",
			props.Dep, props.Dep)
	}

	module, _, _ := adapter_conf.ModuleSpec()

	return addDepAction.AddRemoteDepInternal(deps, io, api.AddDepProps{
		Path:            props.Path,
		Dep:             module + "@" + props.Version,
		As:              props.Dep,
		RemoteAvailable: props.RemoteAvailable,
	})
}
