package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// requireModulePath returns the module-path field of a require entry
// ("github.com/x/y v1.2.3 // indirect" -> "github.com/x/y").
func requireModulePath(deps *deps.Deps, require string) string {
	fields := deps.Stringsdeps.Fields(require)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func addRequire(deps *deps.Deps, conf *ModuleConf, require string) {
	module := requireModulePath(deps, require)
	for i, existing := range conf.Requires {
		if requireModulePath(deps, existing) == module {
			conf.Requires[i] = require
			return
		}
	}
	conf.Requires = append(conf.Requires, require)
}

func removeRequire(deps *deps.Deps, conf *ModuleConf, module string) {
	kept := conf.Requires[:0:0]
	for _, existing := range conf.Requires {
		if requireModulePath(deps, existing) == module {
			continue
		}
		kept = append(kept, existing)
	}
	conf.Requires = kept
}
