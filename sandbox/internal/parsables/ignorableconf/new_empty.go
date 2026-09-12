package ignorableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *IgnorableConf {
	items := &IgnorableConf{
		Paths: make([]string, 0),
	}
	BindMethods(sandbox, items)
	return items
}
