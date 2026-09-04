package docpropsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func BindMethods(deps *deps.Deps, doc_props_conf *DocPropsConf) {

	doc_props_conf.AddTheme = func(id string) {
		for _, theme := range doc_props_conf.Themes {
			if theme == id {
				return
			}
		}
		doc_props_conf.Themes = append(doc_props_conf.Themes, id)
	}

	doc_props_conf.Render = func() string {
		return Render(deps, doc_props_conf)
	}
}
