package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Render(deps *deps.Deps, project_conf *ProjectConf) string {
	obj := deps.SerializeLib.CreateObject()
	obj.AddItemToObject("name", project_conf.Name)
	obj.AddItemToObject("version", project_conf.Version)
	obj.AddItemToObject("description", project_conf.Description)

	return deps.SerializeLib.SerializeToYaml(obj)
}
