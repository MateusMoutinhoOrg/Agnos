package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func BindMethods(deps *deps.Deps, conf *CommandConf) {
	conf.Render = func() string {
		return Render(deps, conf)
	}
}
