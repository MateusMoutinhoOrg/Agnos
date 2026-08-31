package ignorableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Render(deps *deps.Deps, items *IgnorableConf) string {
	arr := deps.SerializeLib.CreateArray()

	for _, p := range items.Paths {
		arr.AddItemToArray(deps.SerializeLib.CreateString(p))
	}

	return deps.SerializeLib.SerializeToYaml(arr)
}
