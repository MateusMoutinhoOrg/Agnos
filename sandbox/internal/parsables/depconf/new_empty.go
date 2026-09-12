package depconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *DepConf {
	dep_conf := &DepConf{}
	BindMethods(sandbox, dep_conf)
	return dep_conf
}
