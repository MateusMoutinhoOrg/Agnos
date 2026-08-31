package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *ProjectConf {

	project_conf := &ProjectConf{
		Version: "v0.0.0",
	}
	BindMethods(deps, project_conf)
	return project_conf
}
