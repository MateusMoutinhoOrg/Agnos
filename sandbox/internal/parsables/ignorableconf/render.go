package ignorableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Render(deps *deps.Deps, items *IgnorableConf) string {
	arr := deps.Serializables.CreateArray()

	for _, p := range items.Paths {
		arr.AddItemToArray(deps.Serializables.CreateString(p))
	}

	return deps.Serializables.SerializeToYaml(arr)
}
