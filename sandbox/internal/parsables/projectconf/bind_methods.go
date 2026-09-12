package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, project_conf *ProjectConf) {
	project_conf.Render = func() string {
		return Render(sandbox, project_conf)
	}
}
