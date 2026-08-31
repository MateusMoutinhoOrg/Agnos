package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func BindMethods(deps *deps.Deps, conf *ModuleConf) {
	conf.Render = func() string {
		return Render(deps, conf)
	}
}
