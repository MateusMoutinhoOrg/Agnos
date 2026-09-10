package remove_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_adapter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveDepInternal uninstalls one dep: the contract, and the adapters that
// fill it.
//
// A contract may now have more than one implementation, so a dep that still
// has adapters installed is refused with those adapters named — the same
// answer `remove-route` gives for a route that has a page. WithAdapters is the
// caller saying to take them all; removing only one is `remove-adapter`.
func RemoveDepInternal(deps *deps.Deps, io *smartio.SmartIO, props api.RemoveDepProps) error {
	deps.Std.Log("remove-dep started with path %s dep %s \n", props.Path, props.Dep)

	contract := utils.ContractsDir + "/" + props.Dep
	if !io.IsDir(contract) {
		return deps.Std.Errorf("dep %q is not installed", props.Dep)
	}

	adapters := utils.AdaptersFillingDep(deps, io, props.Dep)

	if len(adapters) > 0 && !props.WithAdapters {
		return deps.Std.Errorf("dep %q is filled by %s: remove one with `agnos remove-adapter <adapter>`, or all of them with `agnos remove-dep %s --with-adapters`",
			props.Dep, deps.Stringsdeps.Join(adapters, ", "), props.Dep)
	}

	for _, adapter := range adapters {
		adapter_conf, err := utils.LoadAdapterConf(deps, io, adapter)
		if err != nil {
			return err
		}
		if err := removeAdapterAction.Uninstall(deps, io, adapter_conf); err != nil {
			return err
		}
	}

	utils.RemoveTree(deps, io, []string{contract})
	return nil
}
