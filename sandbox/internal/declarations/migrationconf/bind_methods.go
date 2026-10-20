package migrationconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *MigrationConf) {
	conf.Render = func() string {
		return Render(sandbox, conf)
	}
}
