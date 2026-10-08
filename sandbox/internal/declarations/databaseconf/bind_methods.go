package databaseconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *DatabaseConf) {
	conf.Render = func() string {
		return Render(sandbox, conf)
	}
}
