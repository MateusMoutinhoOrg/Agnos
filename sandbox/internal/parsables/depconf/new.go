package depconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func New(deps *deps.Deps, content string) (*DepConf, error) {

	if content == "" {
		return nil, deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	dep_specs, parse_error := deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !dep_specs.IsObject() {
		return nil, deps.Std.Errorf("dep_specs is not an object")
	}

	dep_conf := &DepConf{}

	for _, field := range []struct {
		key    string
		target *string
	}{
		{"name", &dep_conf.Name},
		{"field", &dep_conf.Field},
		{"help", &dep_conf.Help},
		{"default-adapter", &dep_conf.DefaultAdapter},
	} {
		item, _ := dep_specs.GetObjectItem(field.key)
		if item == nil || item.IsNull() {
			continue
		}
		value, err := item.GetString()
		if err != nil {
			return nil, deps.Std.Errorf("%s is not a string", field.key)
		}
		*field.target = value
	}

	BindMethods(deps, dep_conf)
	return dep_conf, nil
}
