package userconfig

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/serializebles"
)

type Theme struct {
	Name        string
	Description string
}

type ThemesConf struct {
	Themes []Theme

	AddTheme func(name string, description string) error
	GetTheme func(name string) (*Theme, error)
	Persist  func() error
}

func NewThemesConf(sandbox *api.SandBox, path string) (*ThemesConf, error) {

	var themes_specs *serializibles.SerializibleObject
	if sandbox.Deps.IoLib.IsFile(path) {
		content_bytes, fileerror := sandbox.Deps.IoLib.ReadFile(path)
		if fileerror != nil {
			return nil, fileerror
		}
		specs, parse_error := sandbox.Deps.SerializeLib.ParseYaml(string(content_bytes))
		if parse_error != nil {
			return nil, parse_error
		}
		themes_specs = specs

	} else {
		themes_specs = sandbox.Deps.SerializeLib.CreateArray()
	}

	if !themes_specs.IsArray() {
		return nil, sandbox.Deps.Errorf("themes_specs is not an array")
	}

	themes_conf := ThemesConf{
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

		themes_conf.Themes = append(themes_conf.Themes, theme)
	}

	themes_conf.GetTheme = func(name string) (*Theme, error) {
		for i, theme := range themes_conf.Themes {
			if theme.Name == name {
				return &themes_conf.Themes[i], nil
			}
		}
		return nil, sandbox.Deps.Errorf("theme not found")
	}

	themes_conf.AddTheme = func(name string, description string) error {
		_, err := themes_conf.GetTheme(name)
		if err == nil {
			return sandbox.Deps.Errorf("theme already exists")
		}

		themes_conf.Themes = append(themes_conf.Themes, Theme{
			Name:        name,
			Description: description,
		})

		return nil
	}

	themes_conf.Persist = func() error {
		new_themes_specs := sandbox.Deps.SerializeLib.CreateArray()

		for _, theme := range themes_conf.Themes {
			theme_obj := sandbox.Deps.SerializeLib.CreateObject()
			theme_obj.AddItemToObject("name", theme.Name)
			theme_obj.AddItemToObject("description", theme.Description)

			new_themes_specs.AddItemToArray(theme_obj)
		}

		bytes := sandbox.Deps.SerializeLib.SerializeToYaml(new_themes_specs)
		return sandbox.Deps.IoLib.WriteFile(path, []byte(bytes))
	}

	return &themes_conf, nil
}
