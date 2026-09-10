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

	for index, directive := range conf.Directives {
		if index == 0 {
			builder = separate(deps, builder)
		}
		builder += directive + "\n"
	}

	return builder
}

// separate leaves exactly one blank line between what is written and the
// directives that follow, however many the sections above happened to end with.
func separate(deps *deps.Deps, builder string) string {
	for deps.Stringsdeps.HasSuffix(builder, "\n") {
		builder = builder[:len(builder)-1]
	}
	return builder + "\n\n"
}
