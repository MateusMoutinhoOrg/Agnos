package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *CommandConf) {
	conf.Render = func() string {
		return Render(sandbox, conf)
	}
}
