package extensionsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *ExtensionsConf) {

	conf.SetEnabled = func(name string, enabled bool) {
		for i, extension := range conf.Extensions {
			if extension.Name == name {
				conf.Extensions[i].Enabled = enabled
				return
			}
		}
		conf.Extensions = append(conf.Extensions, Extension{
			Name:    name,
			Enabled: enabled,
		})
	}

	conf.IsEnabled = func(name string) bool {
		for _, extension := range conf.Extensions {
			if extension.Name == name {
				return extension.Enabled
			}
		}
		return false
	}

	conf.Has = func(name string) bool {
		for _, extension := range conf.Extensions {
			if extension.Name == name {
				return true
			}
		}
		return false
	}

	conf.Render = func() string {
		return Render(sandbox, conf)
	}
}
