package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, project_conf *ProjectConf) string {
	obj := sandbox.Deps.SerializableDeps.CreateObject()
	obj.AddItemToObject("project-name", project_conf.ProjectName)
	if project_conf.Version == "" {
		obj.AddItemToObject("version", nil)
	} else {
		obj.AddItemToObject("version", project_conf.Version)
	}
	return sandbox.Deps.SerializableDeps.SerializeToYaml(obj)
}
