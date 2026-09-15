package depconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, dep_conf *DepConf) {
	dep_conf.Render = func() string {
		return Render(sandbox, dep_conf)
	}
}
