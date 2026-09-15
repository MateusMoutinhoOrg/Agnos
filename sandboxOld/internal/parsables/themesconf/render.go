package themesconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, themes_conf *ThemesConf) string {
	new_themes_specs := sandbox.Deps.Serializables.CreateArray()

	for _, theme := range themes_conf.Themes {
		theme_obj := sandbox.Deps.Serializables.CreateObject()
		theme_obj.AddItemToObject("name", theme.Name)
		theme_obj.AddItemToObject("id", theme.Id)
		theme_obj.AddItemToObject("description", theme.Description)

		new_themes_specs.AddItemToArray(theme_obj)
	}

	return sandbox.Deps.Serializables.SerializeToYaml(new_themes_specs)
}
