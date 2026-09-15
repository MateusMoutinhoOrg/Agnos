package docpropsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, doc_props_conf *DocPropsConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()
	obj.AddItemToObject("name", doc_props_conf.Name)
	obj.AddItemToObject("description", doc_props_conf.Description)

	if len(doc_props_conf.Themes) > 0 {
		themes := sandbox.Deps.Serializables.CreateArray()
		for _, theme := range doc_props_conf.Themes {
			themes.AddItemToArray(theme)
		}
		obj.AddItemToObject("themes", themes)
	}

	if doc_props_conf.HasOrder {
		obj.AddItemToObject("order", int64(doc_props_conf.Order))
	}

	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}
