package adapterconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func New(sandbox *api.Sandbox, content string) (*AdapterConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	adapter_specs, parse_error := sandbox.Deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !adapter_specs.IsObject() {
		return nil, sandbox.Deps.Std.Errorf("adapter_specs is not an object")
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
			return nil, sandbox.Deps.Std.Errorf("%s is not a string", field.key)
		}
		*field.target = value
	}

	BindMethods(sandbox, adapter_conf)
	return adapter_conf, nil
}
