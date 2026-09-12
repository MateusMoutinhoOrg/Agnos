package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// requireModulePath returns the module-path field of a require entry
// ("github.com/x/y v1.2.3 // indirect" -> "github.com/x/y").
func requireModulePath(sandbox *api.Sandbox, require string) string {
	fields := sandbox.Deps.Stringsdeps.Fields(require)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func addRequire(sandbox *api.Sandbox, conf *ModuleConf, require string) {
	module := requireModulePath(sandbox, require)
	for i, existing := range conf.Requires {
		if requireModulePath(sandbox, existing) == module {
			conf.Requires[i] = require
			return
		}
	}
	conf.Requires = append(conf.Requires, require)
}

func removeRequire(sandbox *api.Sandbox, conf *ModuleConf, module string) {
	kept := conf.Requires[:0:0]
	for _, existing := range conf.Requires {
		if requireModulePath(sandbox, existing) == module {
			continue
		}
		kept = append(kept, existing)
	}
	conf.Requires = kept
}
