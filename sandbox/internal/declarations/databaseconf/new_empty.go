package databaseconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *DatabaseConf {
	conf := &DatabaseConf{
		Version: 1,
		Tables:  []Table{},
	}
	BindMethods(sandbox, conf)
	return conf
}
