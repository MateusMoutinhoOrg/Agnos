package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func Render(deps *deps.Deps, conf *ModuleConf) string {
	builder := ""
	if conf.Module != "" {
		builder += "module " + conf.Module + "\n\n"
	}
	if conf.GoVersion != "" {
		builder += "go " + conf.GoVersion + "\n\n"
	}

	if len(conf.Requires) > 0 {
		builder += "require (\n"
		for _, req := range conf.Requires {
			builder += "\t" + req + "\n"
		}
		builder += ")\n"
	}

	return builder
}
