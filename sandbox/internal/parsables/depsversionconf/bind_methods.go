package depsversionconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func BindMethods(deps *deps.Deps, conf *DepsVersionConf) {
	conf.Get = func(dep string) (string, string, bool) {
		spec, ok := conf.Deps[dep]
		if !ok {
			return "", "", false
		}
		at := deps.Stringsdeps.LastIndex(spec, "@")
		if at < 0 {
			return spec, "", true
		}
		return spec[:at], spec[at+1:], true
	}

	conf.Render = func() string {
		return Render(deps, conf)
	}
}
