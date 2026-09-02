package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *ProjectConf {

	// Version is left empty: a freshly started project has no version yet.
	project_conf := &ProjectConf{}
	BindMethods(deps, project_conf)
	return project_conf
}
