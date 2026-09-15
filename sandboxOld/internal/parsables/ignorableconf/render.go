package ignorableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, items *IgnorableConf) string {
	arr := sandbox.Deps.Serializables.CreateArray()

	for _, p := range items.Paths {
		arr.AddItemToArray(sandbox.Deps.Serializables.CreateString(p))
	}

	return sandbox.Deps.Serializables.SerializeToYaml(arr)
}
