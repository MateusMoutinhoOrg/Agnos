package depsversionconf

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func BindMethods(deps *deps.Deps, conf *DepsVersionConf) {
	conf.Get = func(dep string) (string, string, bool) {
		spec, ok := conf.Deps[dep]
		if !ok {
			return "", "", false
		}
		at := strings.LastIndex(spec, "@")
		if at < 0 {
			return spec, "", true
		}
		return spec[:at], spec[at+1:], true
	}

	conf.Render = func() string {
		return Render(deps, conf)
	}
}
