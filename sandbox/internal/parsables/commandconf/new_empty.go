package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *CommandConf {
	conf := &CommandConf{
		Identifiers: []string{},
		Examples:    []string{},
		Flags:       []Field{},
		Args:        []Field{},
	}
	BindMethods(deps, conf)
	return conf
}
