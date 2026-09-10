package availableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *AvailableConf {
	available_conf := &AvailableConf{Adapters: []string{}}
	BindMethods(deps, available_conf)
	return available_conf
}
