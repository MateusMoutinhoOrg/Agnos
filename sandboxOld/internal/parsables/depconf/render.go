package depconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, dep_conf *DepConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()
	obj.AddItemToObject("name", dep_conf.Name)
	obj.AddItemToObject("field", dep_conf.Field)
	obj.AddItemToObject("help", dep_conf.Help)
	obj.AddItemToObject("default-adapter", dep_conf.DefaultAdapter)
	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}
