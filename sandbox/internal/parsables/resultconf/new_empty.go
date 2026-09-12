package resultconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *ResultConf {
	conf := &ResultConf{}
	BindMethods(sandbox, conf)
	return conf
}
