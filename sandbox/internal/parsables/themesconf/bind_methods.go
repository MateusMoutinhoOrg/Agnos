package themesconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func BindMethods(deps *deps.Deps, themes_conf *ThemesConf) {

	themes_conf.GetTheme = func(name string) (*Theme, error) {
		for i, theme := range themes_conf.Themes {
			if theme.Name == name {
				return &themes_conf.Themes[i], nil
			}
		}
		return nil, deps.Std.Errorf("theme not found")
	}

	themes_conf.AddTheme = func(name string, id string, description string) error {
		_, err := themes_conf.GetTheme(name)
		if err == nil {
			return deps.Std.Errorf("theme already exists")
		}

		themes_conf.Themes = append(themes_conf.Themes, Theme{
			Name:        name,
			Id:          id,
			Description: description,
		})

		return nil
	}

	themes_conf.Render = func() string {
		return Render(deps, themes_conf)
	}
}
