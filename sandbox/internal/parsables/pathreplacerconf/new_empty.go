package pathreplacerconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *PathReplacerConf {
	conf := &PathReplacerConf{
		Entries: make([]PathReplacerEntry, 0),
	}
	BindMethods(deps, conf)
	return conf
}
