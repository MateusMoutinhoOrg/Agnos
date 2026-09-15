package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, conf *ModuleConf) string {
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
			builder = separate(sandbox, builder)
		}
		builder += directive + "\n"
	}

	return builder
}

// separate leaves exactly one blank line between what is written and the
// directives that follow, however many the sections above happened to end with.
func separate(sandbox *api.Sandbox, builder string) string {
	for sandbox.Deps.Stringsdeps.HasSuffix(builder, "\n") {
		builder = builder[:len(builder)-1]
	}
	return builder + "\n\n"
}
