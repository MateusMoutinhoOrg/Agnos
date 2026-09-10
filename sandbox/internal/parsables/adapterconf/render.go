package adapterconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func Render(deps *deps.Deps, adapter_conf *AdapterConf) string {
	obj := deps.Serializables.CreateObject()
	obj.AddItemToObject("name", adapter_conf.Name)
	obj.AddItemToObject("dep", adapter_conf.Dep)
	obj.AddItemToObject("help", adapter_conf.Help)
	obj.AddItemToObject("module", adapter_conf.Module)
	obj.AddItemToObject("origin", adapter_conf.Origin)
	return deps.Serializables.SerializeToYaml(obj)
}
