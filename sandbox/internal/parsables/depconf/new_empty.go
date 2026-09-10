package depconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *DepConf {
	dep_conf := &DepConf{}
	BindMethods(deps, dep_conf)
	return dep_conf
}
