package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func BindMethods(deps *deps.Deps, project_conf *ProjectConf) {
	project_conf.Render = func() string {
		return Render(deps, project_conf)
	}
}
