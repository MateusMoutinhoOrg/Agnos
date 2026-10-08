package set_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetAdapterInternal changes which installed adapter fills one dep's field in
// one binding. It is the only editor of that choice, the way `set-route` is
// the only editor of a route's method: the selection lives in one file, and
// switching it moves every generated bind at once instead of leaving two.
func SetAdapterInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.SetAdapterProps) error {
	binding := props.Binding
	if binding == "" {
		binding = utils.StandardBinding
	}

	sandbox.Deps.StdDeps.Logf("set-adapter started with path %s dep %s adapter %s binding %s \n",
		props.Path, props.Dep, props.Adapter, binding)

	adapter_conf, err := utils.LoadAdapterConf(sandbox, io, props.Adapter)
	if err != nil {
		return sandbox.Deps.StdDeps.Errorf("adapter %q is not installed (run `agnos add-adapter %s`)", props.Adapter, props.Adapter)
	}

	if adapter_conf.Dep != props.Dep {
		return sandbox.Deps.StdDeps.Errorf("adapter %q fills dep %q, not %q", props.Adapter, adapter_conf.Dep, props.Dep)
	}

	if _, err := utils.LoadBindingConf(sandbox, io, binding); err != nil {
		return err
	}

	return utils.SelectAdapter(sandbox, io, binding, props.Adapter)
}
