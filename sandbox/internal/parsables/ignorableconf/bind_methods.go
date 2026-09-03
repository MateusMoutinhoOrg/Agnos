package ignorableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func BindMethods(deps *deps.Deps, items *IgnorableConf) {

	items.AddPath = func(path string) {
		items.Paths = append(items.Paths, path)
	}

	items.IsIgnorable = func(path string) bool {
		for _, p := range items.Paths {
			if p == path {
				return true
			}
		}
		return false
	}

	items.Render = func() string {
		return Render(deps, items)
	}
}
