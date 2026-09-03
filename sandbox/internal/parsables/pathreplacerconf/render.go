package pathreplacerconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Render(deps *deps.Deps, conf *PathReplacerConf) string {
	obj := deps.Serializables.CreateObject()

	for _, entry := range conf.Entries {
		obj.AddItemToObject(entry.Original, entry.Replacement)
	}

	return deps.Serializables.SerializeToYaml(obj)
}
