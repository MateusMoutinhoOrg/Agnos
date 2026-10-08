package bindingconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, binding_conf *BindingConf) {
	binding_conf.Has = func(adapter string) bool {
		for _, candidate := range binding_conf.Adapters {
			if candidate == adapter {
				return true
			}
		}
		return false
	}

	binding_conf.Add = func(adapter string) bool {
		if binding_conf.Has(adapter) {
			return false
		}
		binding_conf.Adapters = append(binding_conf.Adapters, adapter)
		sandbox.Deps.SortDeps.Strings(binding_conf.Adapters)
		return true
	}

	binding_conf.Remove = func(adapter string) bool {
		kept := []string{}
		for _, candidate := range binding_conf.Adapters {
			if candidate != adapter {
				kept = append(kept, candidate)
			}
		}
		if len(kept) == len(binding_conf.Adapters) {
			return false
		}
		binding_conf.Adapters = kept
		return true
	}

	binding_conf.Render = func() string {
		return Render(sandbox, binding_conf)
	}
}
