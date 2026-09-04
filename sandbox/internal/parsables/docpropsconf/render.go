package docpropsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func Render(deps *deps.Deps, doc_props_conf *DocPropsConf) string {
	obj := deps.Serializables.CreateObject()
	obj.AddItemToObject("name", doc_props_conf.Name)
	obj.AddItemToObject("description", doc_props_conf.Description)

	if len(doc_props_conf.Themes) > 0 {
		themes := deps.Serializables.CreateArray()
		for _, theme := range doc_props_conf.Themes {
			themes.AddItemToArray(theme)
		}
		obj.AddItemToObject("themes", themes)
	}

	if doc_props_conf.HasOrder {
		obj.AddItemToObject("order", int64(doc_props_conf.Order))
	}

	return deps.Serializables.SerializeToYaml(obj)
}
