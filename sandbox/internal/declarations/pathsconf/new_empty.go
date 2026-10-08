package pathreplacerconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *PathReplacerConf {
	conf := &PathReplacerConf{
		Entries: make([]PathReplacerEntry, 0),
	}
	BindMethods(sandbox, conf)
	return conf
}
