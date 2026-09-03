package depsversionconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *DepsVersionConf {
	conf := &DepsVersionConf{
		Deps: map[string]string{},
	}
	BindMethods(deps, conf)
	return conf
}
