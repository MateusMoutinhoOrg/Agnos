package availableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, available_conf *AvailableConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()

	adapters := sandbox.Deps.Serializables.CreateArray()
	for _, adapter := range available_conf.Adapters {
		adapters.AddItemToArray(adapter)
	}

	obj.AddItemToObject("adapters", adapters)
	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}
