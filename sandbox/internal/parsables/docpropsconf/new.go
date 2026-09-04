package docpropsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func New(deps *deps.Deps, content string) (*DocPropsConf, error) {

	if content == "" {
		return nil, deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}
	props_specs := specs

	if !props_specs.IsObject() {
		return nil, deps.Std.Errorf("props_specs is not an object")
	}

	doc_props_conf := &DocPropsConf{
		Themes: make([]string, 0),
	}
	var err error

	name_item, _ := props_specs.GetObjectItem("name")
	description_item, _ := props_specs.GetObjectItem("description")
	themes_item, _ := props_specs.GetObjectItem("themes")
	order_item, _ := props_specs.GetObjectItem("order")

	if name_item != nil && !name_item.IsNull() {
		doc_props_conf.Name, err = name_item.GetString()
		if err != nil {
			return nil, deps.Std.Errorf("name is not a string")
		}
	}

	if description_item != nil && !description_item.IsNull() {
		doc_props_conf.Description, err = description_item.GetString()
		if err != nil {
			return nil, deps.Std.Errorf("description is not a string")
		}
	}

	if themes_item != nil && !themes_item.IsNull() {
		if !themes_item.IsArray() {
			return nil, deps.Std.Errorf("themes is not an array")
		}

		size, err := themes_item.GetArraySize()
		if err != nil {
			return nil, deps.Std.Errorf("could not get themes array size")
		}

		for i := 0; i < size; i++ {
			item := themes_item.GetArrayItem(i)
			if item == nil || item.IsNull() {
				continue
			}

			theme, err := item.GetString()
			if err != nil {
				return nil, deps.Std.Errorf("theme is not a string")
			}
			doc_props_conf.Themes = append(doc_props_conf.Themes, theme)
		}
	}

	// order is optional: absent means "list this doc after every ordered one".
	if order_item != nil && !order_item.IsNull() {
		order, err := order_item.GetInt()
		if err != nil {
			return nil, deps.Std.Errorf("order is not an int")
		}
		doc_props_conf.Order = int(order)
		doc_props_conf.HasOrder = true
	}

	BindMethods(deps, doc_props_conf)
	return doc_props_conf, nil
}
