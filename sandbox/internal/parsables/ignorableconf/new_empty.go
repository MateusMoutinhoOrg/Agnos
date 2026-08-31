package ignorableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *IgnorableConf {
	items := &IgnorableConf{
		Paths: make([]string, 0),
	}
	BindMethods(deps, items)
	return items
}
