package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func Render(deps *deps.Deps, project_conf *ProjectConf) string {
	obj := deps.Serializables.CreateObject()
	obj.AddItemToObject("name", project_conf.Name)
	if project_conf.Version == "" {
		obj.AddItemToObject("version", nil)
	} else {
		obj.AddItemToObject("version", project_conf.Version)
	}
	obj.AddItemToObject("description", project_conf.Description)

	return deps.Serializables.SerializeToYaml(obj)
}
