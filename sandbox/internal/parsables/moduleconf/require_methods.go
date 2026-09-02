package moduleconf

import "strings"

// requireModulePath returns the module-path field of a require entry
// ("github.com/x/y v1.2.3 // indirect" -> "github.com/x/y").
func requireModulePath(require string) string {
	fields := strings.Fields(require)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func addRequire(conf *ModuleConf, require string) {
	module := requireModulePath(require)
	for i, existing := range conf.Requires {
		if requireModulePath(existing) == module {
			conf.Requires[i] = require
			return
		}
	}
	conf.Requires = append(conf.Requires, require)
}

func removeRequire(conf *ModuleConf, module string) {
	kept := conf.Requires[:0:0]
	for _, existing := range conf.Requires {
		if requireModulePath(existing) == module {
			continue
		}
		kept = append(kept, existing)
	}
	conf.Requires = kept
}
