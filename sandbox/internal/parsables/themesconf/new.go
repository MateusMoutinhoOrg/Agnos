package themesconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func New(deps *deps.Deps, content string) (*ThemesConf, error) {

	if content == "" {
		return nil, deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := deps.SerializeLib.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}
	themes_specs := specs

	if !themes_specs.IsArray() {
		return nil, deps.Std.Errorf("themes_specs is not an array")
	}

	themes_conf := &ThemesConf{
		Themes: make([]Theme, 0),
	}

	size, err := themes_specs.GetArraySize()
	if err != nil {
		return nil, deps.Std.Errorf("could not get themes array size")
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
				return nil, deps.Std.Errorf("name is not a string")
			}
		}

		if description_item != nil && !description_item.IsNull() {
			theme.Description, err = description_item.GetString()
			if err != nil {
				return nil, deps.Std.Errorf("description is not a string")
			}
		}

		if id_item != nil && !id_item.IsNull() {
			theme.Id, err = id_item.GetString()
			if err != nil {
				return nil, deps.Std.Errorf("id is not a string")
			}
		}

		themes_conf.Themes = append(themes_conf.Themes, theme)
	}

	BindMethods(deps, themes_conf)
	return themes_conf, nil
}
