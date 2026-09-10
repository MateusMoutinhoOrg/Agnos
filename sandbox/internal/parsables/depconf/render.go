package depconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func Render(deps *deps.Deps, dep_conf *DepConf) string {
	obj := deps.Serializables.CreateObject()
	obj.AddItemToObject("name", dep_conf.Name)
	obj.AddItemToObject("field", dep_conf.Field)
	obj.AddItemToObject("help", dep_conf.Help)
	obj.AddItemToObject("default-adapter", dep_conf.DefaultAdapter)
	return deps.Serializables.SerializeToYaml(obj)
}
