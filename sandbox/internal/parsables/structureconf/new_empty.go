package structureconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *StructureConf {
	structure_conf := &StructureConf{
		Items: make([]Item, 0),
	}
	BindMethods(deps, structure_conf)
	return structure_conf
}
