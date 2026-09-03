package depsversionconf

import (
	"sort"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Render(deps *deps.Deps, conf *DepsVersionConf) string {
	obj := deps.Serializables.CreateObject()

	names := make([]string, 0, len(conf.Deps))
	for name := range conf.Deps {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		obj.AddItemToObject(name, conf.Deps[name])
	}

	return deps.Serializables.SerializeToYaml(obj)
}
