package resultconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *ResultConf {
	conf := &ResultConf{}
	BindMethods(deps, conf)
	return conf
}
