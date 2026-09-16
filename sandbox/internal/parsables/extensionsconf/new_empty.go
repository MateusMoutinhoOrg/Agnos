package extensionsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *ExtensionsConf {
	conf := &ExtensionsConf{
		Extensions: make([]Extension, 0),
	}
	BindMethods(sandbox, conf)
	return conf
}
