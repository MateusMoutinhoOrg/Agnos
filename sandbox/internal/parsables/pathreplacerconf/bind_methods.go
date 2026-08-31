package pathreplacerconf

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func BindMethods(deps *deps.Deps, conf *PathReplacerConf) {

	conf.AddEntry = func(original string, replacement string) {
		conf.Entries = append(conf.Entries, PathReplacerEntry{
			Original:    original,
			Replacement: replacement,
		})
	}

	conf.Format = func(path string) string {
		result := path
		for _, entry := range conf.Entries {
			result = strings.ReplaceAll(result, entry.Original, entry.Replacement)
		}
		return result
	}

	conf.Render = func() string {
		return Render(deps, conf)
	}
}
