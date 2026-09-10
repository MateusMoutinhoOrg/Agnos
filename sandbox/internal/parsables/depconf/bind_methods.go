package depconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func BindMethods(deps *deps.Deps, dep_conf *DepConf) {
	dep_conf.Render = func() string {
		return Render(deps, dep_conf)
	}
}
