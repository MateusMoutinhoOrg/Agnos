package pathsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *PathsConf {
	conf := &PathsConf{
		Entries: make([]PathReplacerEntry, 0),
	}
	BindMethods(sandbox, conf)
	return conf
}
