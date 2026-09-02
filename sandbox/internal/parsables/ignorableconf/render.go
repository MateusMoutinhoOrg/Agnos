package ignorableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Render(deps *deps.Deps, items *IgnorableConf) string {
	arr := deps.Serializebles.CreateArray()

	for _, p := range items.Paths {
		arr.AddItemToArray(deps.Serializebles.CreateString(p))
	}

	return deps.Serializebles.SerializeToYaml(arr)
}
