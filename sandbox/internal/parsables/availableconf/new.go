package availableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func New(sandbox *api.Sandbox, content string) (*AvailableConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	available_specs, parse_error := sandbox.Deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !available_specs.IsObject() {
		return nil, sandbox.Deps.Std.Errorf("available_specs is not an object")
	}

	available_conf := &AvailableConf{Adapters: []string{}}

	adapters_item, _ := available_specs.GetObjectItem("adapters")
	if adapters_item != nil && !adapters_item.IsNull() {
		if !adapters_item.IsArray() {
			return nil, sandbox.Deps.Std.Errorf("adapters is not an array")
		}
		size, err := adapters_item.GetArraySize()
		if err != nil {
			return nil, err
		}
		for index := 0; index < size; index++ {
			item := adapters_item.GetArrayItem(index)
			if item == nil || item.IsNull() {
				continue
			}
			value, err := item.GetString()
			if err != nil {
				return nil, sandbox.Deps.Std.Errorf("adapters entry %d is not a string", index)
			}
			available_conf.Adapters = append(available_conf.Adapters, value)
		}
	}

	BindMethods(sandbox, available_conf)
	return available_conf, nil
}
