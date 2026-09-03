package moduleconf

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func Render(deps *deps.Deps, conf *ModuleConf) string {
	var builder strings.Builder
	if conf.Module != "" {
		builder.WriteString("module " + conf.Module + "\n\n")
	}
	if conf.GoVersion != "" {
		builder.WriteString("go " + conf.GoVersion + "\n\n")
	}

	if len(conf.Requires) > 0 {
		builder.WriteString("require (\n")
		for _, req := range conf.Requires {
			builder.WriteString("\t" + req + "\n")
		}
		builder.WriteString(")\n")
	}

	return builder.String()
}
