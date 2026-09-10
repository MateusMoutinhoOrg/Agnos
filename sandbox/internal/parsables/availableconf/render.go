package availableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func Render(deps *deps.Deps, available_conf *AvailableConf) string {
	obj := deps.Serializables.CreateObject()

	adapters := deps.Serializables.CreateArray()
	for _, adapter := range available_conf.Adapters {
		adapters.AddItemToArray(adapter)
	}

	obj.AddItemToObject("adapters", adapters)
	return deps.Serializables.SerializeToYaml(obj)
}
