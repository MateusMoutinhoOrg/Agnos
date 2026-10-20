package migrationconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *MigrationConf {
	conf := &MigrationConf{
		Ops: []Op{},
	}
	BindMethods(sandbox, conf)
	return conf
}
