package docconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, doc_props_conf *DocConf) string {
	obj := sandbox.Deps.SerializableDeps.CreateObject()
	obj.AddItemToObject("name", doc_props_conf.Name)
	obj.AddItemToObject("description", doc_props_conf.Description)

	if len(doc_props_conf.Themes) > 0 {
		themes := sandbox.Deps.SerializableDeps.CreateArray()
		for _, theme := range doc_props_conf.Themes {
			themes.AddItemToArray(theme)
		}
		obj.AddItemToObject("themes", themes)
	}

	if doc_props_conf.HasOrder {
		obj.AddItemToObject("order", int64(doc_props_conf.Order))
	}

	return sandbox.Deps.SerializableDeps.SerializeToYaml(obj)
}
