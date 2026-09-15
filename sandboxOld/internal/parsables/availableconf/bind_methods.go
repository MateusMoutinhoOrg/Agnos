package availableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, available_conf *AvailableConf) {
	available_conf.Has = func(adapter string) bool {
		for _, candidate := range available_conf.Adapters {
			if candidate == adapter {
				return true
			}
		}
		return false
	}

	available_conf.Add = func(adapter string) bool {
		if available_conf.Has(adapter) {
			return false
		}
		available_conf.Adapters = append(available_conf.Adapters, adapter)
		sandbox.Deps.Sortdeps.Strings(available_conf.Adapters)
		return true
	}

	available_conf.Remove = func(adapter string) bool {
		kept := []string{}
		for _, candidate := range available_conf.Adapters {
			if candidate != adapter {
				kept = append(kept, candidate)
			}
		}
		if len(kept) == len(available_conf.Adapters) {
			return false
		}
		available_conf.Adapters = kept
		return true
	}

	available_conf.Render = func() string {
		return Render(sandbox, available_conf)
	}
}
