package pathreplacerconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *PathReplacerConf {
	conf := &PathReplacerConf{
		Entries: make([]PathReplacerEntry, 0),
	}
	BindMethods(deps, conf)
	return conf
}
