package bindingconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func New(sandbox *api.Sandbox, content string) (*BindingConf, error) {

	if content == "" {
		return nil, sandbox.Deps.StdDeps.Errorf("content cannot be empty, use NewEmpty instead")
	}

	binding_specs, parse_error := sandbox.Deps.SerializableDeps.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !binding_specs.IsObject() {
		return nil, sandbox.Deps.StdDeps.Errorf("binding_specs is not an object")
	}

	binding_conf := &BindingConf{Adapters: []string{}}

	adapters_item, _ := binding_specs.GetObjectItem("adapters")
	if adapters_item != nil && !adapters_item.IsNull() {
		if !adapters_item.IsArray() {
			return nil, sandbox.Deps.StdDeps.Errorf("adapters is not an array")
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
				return nil, sandbox.Deps.StdDeps.Errorf("adapters entry %d is not a string", index)
			}
			binding_conf.Adapters = append(binding_conf.Adapters, value)
		}
	}

	BindMethods(sandbox, binding_conf)
	return binding_conf, nil
}
