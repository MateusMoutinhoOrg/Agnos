package themesconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Render(deps *deps.Deps, themes_conf *ThemesConf) string {
	new_themes_specs := deps.Serializebles.CreateArray()

	for _, theme := range themes_conf.Themes {
		theme_obj := deps.Serializebles.CreateObject()
		theme_obj.AddItemToObject("name", theme.Name)
		theme_obj.AddItemToObject("id", theme.Id)
		theme_obj.AddItemToObject("description", theme.Description)

		new_themes_specs.AddItemToArray(theme_obj)
	}

	return deps.Serializebles.SerializeToYaml(new_themes_specs)
}
