package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *CommandConf {
	conf := &CommandConf{
		Identifiers: []string{},
		Examples:    []string{},
		Flags:       []Field{},
		Args:        []Field{},
	}
	BindMethods(sandbox, conf)
	return conf
}
