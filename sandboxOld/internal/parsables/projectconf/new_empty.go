package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *ProjectConf {

	// Version is left empty: a freshly started project has no version yet.
	project_conf := &ProjectConf{}
	BindMethods(sandbox, project_conf)
	return project_conf
}
