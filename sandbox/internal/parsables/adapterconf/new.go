package adapterconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func New(deps *deps.Deps, content string) (*AdapterConf, error) {

	if content == "" {
		return nil, deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	adapter_specs, parse_error := deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !adapter_specs.IsObject() {
		return nil, deps.Std.Errorf("adapter_specs is not an object")
	}

	adapter_conf := &AdapterConf{Origin: OriginCatalog}

	for _, field := range []struct {
		key    string
		target *string
	}{
		{"name", &adapter_conf.Name},
		{"dep", &adapter_conf.Dep},
		{"help", &adapter_conf.Help},
		{"module", &adapter_conf.Module},
		{"origin", &adapter_conf.Origin},
	} {
		item, _ := adapter_specs.GetObjectItem(field.key)
		if item == nil || item.IsNull() {
			continue
		}
		value, err := item.GetString()
		if err != nil {
			return nil, deps.Std.Errorf("%s is not a string", field.key)
		}
		*field.target = value
	}

	BindMethods(deps, adapter_conf)
	return adapter_conf, nil
}
