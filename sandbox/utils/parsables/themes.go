package parsables

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

type Theme struct {
	Name        string
	Id          string
	Description string
}

type ThemesConf struct {
	Themes []Theme

	AddTheme func(name string, id string, description string) error
	GetTheme func(name string) (*Theme, error)
	Render   func() string
}

func addThemesConfMethods(sandbox *lib.SandBox, themes_conf *ThemesConf) {

	themes_conf.GetTheme = func(name string) (*Theme, error) {
		for i, theme := range themes_conf.Themes {
			if theme.Name == name {
				return &themes_conf.Themes[i], nil
			}
		}
		return nil, sandbox.Deps.Errorf("theme not found")
	}

	themes_conf.AddTheme = func(name string, id string, description string) error {
		_, err := themes_conf.GetTheme(name)
		if err == nil {
			return sandbox.Deps.Errorf("theme already exists")
		}

		themes_conf.Themes = append(themes_conf.Themes, Theme{
			Name:        name,
			Id:          id,
			Description: description,
		})

		return nil
	}

	themes_conf.Render = func() string {
		new_themes_specs := sandbox.Deps.SerializeLib.CreateArray()

		for _, theme := range themes_conf.Themes {
			theme_obj := sandbox.Deps.SerializeLib.CreateObject()
			theme_obj.AddItemToObject("name", theme.Name)
			theme_obj.AddItemToObject("id", theme.Id)
			theme_obj.AddItemToObject("description", theme.Description)

			new_themes_specs.AddItemToArray(theme_obj)
		}

		return sandbox.Deps.SerializeLib.SerializeToYaml(new_themes_specs)
	}
}

func NewThemesConf(sandbox *lib.SandBox, content string) (*ThemesConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Errorf("content cannot be empty, use NewThemesConfEmpty instead")
	}

	specs, parse_error := sandbox.Deps.SerializeLib.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}
	themes_specs := specs

	if !themes_specs.IsArray() {
		return nil, sandbox.Deps.Errorf("themes_specs is not an array")
	}

	themes_conf := &ThemesConf{
		Themes: make([]Theme, 0),
	}

	size, err := themes_specs.GetArraySize()
	if err != nil {
		return nil, sandbox.Deps.Errorf("could not get themes array size")
	}

	for i := 0; i < size; i++ {
		item := themes_specs.GetArrayItem(i)
		if item == nil || !item.IsObject() {
			continue
		}

		theme := Theme{}
		name_item, _ := item.GetObjectItem("name")
		id_item, _ := item.GetObjectItem("id")
		description_item, _ := item.GetObjectItem("description")

		if name_item != nil && !name_item.IsNull() {
			theme.Name, err = name_item.GetString()
			if err != nil {
				return nil, sandbox.Deps.Errorf("name is not a string")
			}
		}

		if description_item != nil && !description_item.IsNull() {
			theme.Description, err = description_item.GetString()
			if err != nil {
				return nil, sandbox.Deps.Errorf("description is not a string")
			}
		}

		if id_item != nil && !id_item.IsNull() {
			theme.Id, err = id_item.GetString()
			if err != nil {
				return nil, sandbox.Deps.Errorf("id is not a string")
			}
		}

		themes_conf.Themes = append(themes_conf.Themes, theme)
	}

	addThemesConfMethods(sandbox, themes_conf)
	return themes_conf, nil
}

func NewThemesConfEmpty(sandbox *lib.SandBox) *ThemesConf {
	themes_conf := &ThemesConf{
		Themes: make([]Theme, 0),
	}
	addThemesConfMethods(sandbox, themes_conf)
	return themes_conf
}
