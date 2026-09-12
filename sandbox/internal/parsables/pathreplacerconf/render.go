package pathreplacerconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, conf *PathReplacerConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()

	for _, entry := range conf.Entries {
		obj.AddItemToObject(entry.Original, entry.Replacement)
	}

	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}
