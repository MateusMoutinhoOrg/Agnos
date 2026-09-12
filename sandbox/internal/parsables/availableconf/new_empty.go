package availableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *AvailableConf {
	available_conf := &AvailableConf{Adapters: []string{}}
	BindMethods(sandbox, available_conf)
	return available_conf
}
