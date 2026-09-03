package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *ModuleConf {
	conf := &ModuleConf{
		Requires: []string{},
	}
	BindMethods(deps, conf)
	return conf
}
