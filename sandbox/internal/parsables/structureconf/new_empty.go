package structureconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *StructureConf {
	structure_conf := &StructureConf{
		Items: make([]Item, 0),
	}
	BindMethods(sandbox, structure_conf)
	return structure_conf
}
