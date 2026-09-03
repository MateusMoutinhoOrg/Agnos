package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func BindMethods(deps *deps.Deps, conf *ModuleConf) {
	conf.AddRequire = func(require string) {
		addRequire(conf, require)
	}
	conf.RemoveRequire = func(module string) {
		removeRequire(conf, module)
	}
	conf.Render = func() string {
		return Render(deps, conf)
	}
}
