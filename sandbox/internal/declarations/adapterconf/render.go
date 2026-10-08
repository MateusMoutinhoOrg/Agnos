package adapterconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, adapter_conf *AdapterConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()
	obj.AddItemToObject("name", adapter_conf.Name)
	obj.AddItemToObject("dep", adapter_conf.Dep)
	obj.AddItemToObject("help", adapter_conf.Help)
	obj.AddItemToObject("module", adapter_conf.Module)
	obj.AddItemToObject("origin", adapter_conf.Origin)
	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}
