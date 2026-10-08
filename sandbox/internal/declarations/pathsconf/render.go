package pathsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, conf *PathsConf) string {
	obj := sandbox.Deps.SerializableDeps.CreateObject()

	for _, entry := range conf.Entries {
		obj.AddItemToObject(entry.Original, entry.Replacement)
	}

	return sandbox.Deps.SerializableDeps.SerializeToYaml(obj)
}
