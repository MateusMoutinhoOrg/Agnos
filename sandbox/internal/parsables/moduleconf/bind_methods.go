package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *ModuleConf) {
	conf.AddRequire = func(require string) {
		addRequire(sandbox, conf, require)
	}
	conf.RemoveRequire = func(module string) {
		removeRequire(sandbox, conf, module)
	}
	conf.Render = func() string {
		return Render(sandbox, conf)
	}
}
