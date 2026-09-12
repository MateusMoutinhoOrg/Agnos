package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *ModuleConf {
	conf := &ModuleConf{
		Requires: []string{},
	}
	BindMethods(sandbox, conf)
	return conf
}
