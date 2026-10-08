package bindingconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *BindingConf {
	binding_conf := &BindingConf{Adapters: []string{}}
	BindMethods(sandbox, binding_conf)
	return binding_conf
}
