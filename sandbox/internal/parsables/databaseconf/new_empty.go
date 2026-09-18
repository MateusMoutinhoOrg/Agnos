package databaseconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *DatabaseConf {
	conf := &DatabaseConf{
		Tables: []Table{},
	}
	BindMethods(sandbox, conf)
	return conf
}
