package adapterconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *AdapterConf {
	adapter_conf := &AdapterConf{Origin: OriginCatalog}
	BindMethods(sandbox, adapter_conf)
	return adapter_conf
}
