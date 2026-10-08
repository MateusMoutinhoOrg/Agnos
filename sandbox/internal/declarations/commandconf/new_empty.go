package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *CommandConf {
	conf := &CommandConf{
		Priority:    DefaultPriority,
		HasPriority: true,
		Strict:      true,
		Args:        []Arg{},
		Flags:       []Flag{},
		Examples:    []string{},
		Legacy:      []string{},
	}
	BindMethods(sandbox, conf)
	return conf
}
