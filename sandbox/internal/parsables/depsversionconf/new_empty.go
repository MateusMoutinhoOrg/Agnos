package depsversionconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *DepsVersionConf {
	conf := &DepsVersionConf{
		Deps: map[string]string{},
	}
	BindMethods(deps, conf)
	return conf
}
