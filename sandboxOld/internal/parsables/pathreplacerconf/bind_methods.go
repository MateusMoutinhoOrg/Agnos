package pathreplacerconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *PathReplacerConf) {

	conf.AddEntry = func(original string, replacement string) {
		conf.Entries = append(conf.Entries, PathReplacerEntry{
			Original:    original,
			Replacement: replacement,
		})
	}

	conf.Format = func(path string) string {
		result := path
		for _, entry := range conf.Entries {
			result = sandbox.Deps.Stringsdeps.ReplaceAll(result, entry.Original, entry.Replacement)
		}
		return result
	}

	conf.Render = func() string {
		return Render(sandbox, conf)
	}
}
