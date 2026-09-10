package adapterconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *AdapterConf {
	adapter_conf := &AdapterConf{Origin: OriginCatalog}
	BindMethods(deps, adapter_conf)
	return adapter_conf
}
