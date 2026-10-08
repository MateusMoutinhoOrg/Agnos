package bindingconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, binding_conf *BindingConf) string {
	obj := sandbox.Deps.SerializableDeps.CreateObject()

	adapters := sandbox.Deps.SerializableDeps.CreateArray()
	for _, adapter := range binding_conf.Adapters {
		adapters.AddItemToArray(adapter)
	}

	obj.AddItemToObject("adapters", adapters)
	return sandbox.Deps.SerializableDeps.SerializeToYaml(obj)
}
