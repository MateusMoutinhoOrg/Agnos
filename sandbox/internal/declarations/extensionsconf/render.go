package extensionsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, conf *ExtensionsConf) string {
	obj := sandbox.Deps.SerializableDeps.CreateObject()

	for _, extension := range conf.Extensions {
		obj.AddItemToObject(extension.Name, extension.Enabled)
	}

	return sandbox.Deps.SerializableDeps.SerializeToYaml(obj)
}
