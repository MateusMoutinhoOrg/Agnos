package adapterconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, adapter_conf *AdapterConf) {
	adapter_conf.ModuleSpec = func() (string, string, bool) {
		if adapter_conf.Module == "" {
			return "", "", false
		}
		at := sandbox.Deps.Stringsdeps.LastIndex(adapter_conf.Module, "@")
		if at < 0 {
			return adapter_conf.Module, "", true
		}
		return adapter_conf.Module[:at], adapter_conf.Module[at+1:], true
	}

	adapter_conf.Render = func() string {
		return Render(sandbox, adapter_conf)
	}
}
