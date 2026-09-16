package extensionsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, conf *ExtensionsConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()

	for _, extension := range conf.Extensions {
		obj.AddItemToObject(extension.Name, extension.Enabled)
	}

	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}
