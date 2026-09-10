package set_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetAdapterInternal changes which installed adapter fills one dep's field in
// one available. It is the only editor of that choice, the way `set-route` is
// the only editor of a route's method: the selection lives in one file, and
// switching it moves every generated bind at once instead of leaving two.
func SetAdapterInternal(deps *deps.Deps, io *smartio.SmartIO, props api.SetAdapterProps) error {
	available := props.Available
	if available == "" {
		available = utils.StandardAvailable
	}

	deps.Std.Log("set-adapter started with path %s dep %s adapter %s available %s \n",
		props.Path, props.Dep, props.Adapter, available)

	adapter_conf, err := utils.LoadAdapterConf(deps, io, props.Adapter)
	if err != nil {
		return deps.Std.Errorf("adapter %q is not installed (run `agnos add-adapter %s`)", props.Adapter, props.Adapter)
	}

	if adapter_conf.Dep != props.Dep {
		return deps.Std.Errorf("adapter %q fills dep %q, not %q", props.Adapter, adapter_conf.Dep, props.Dep)
	}

	if _, err := utils.LoadAvailableConf(deps, io, available); err != nil {
		return err
	}

	return utils.SelectAdapter(deps, io, available, props.Adapter)
}
