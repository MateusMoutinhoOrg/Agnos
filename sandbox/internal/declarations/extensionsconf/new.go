package extensionsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func New(sandbox *api.Sandbox, content string) (*ExtensionsConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := sandbox.Deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsObject() {
		return nil, sandbox.Deps.Std.Errorf("extensions config is not an object")
	}

	conf := &ExtensionsConf{
		Extensions: make([]Extension, 0),
	}

	keys, keys_err := specs.GetKeys()
	if keys_err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not get extensions keys")
	}

	for _, key := range keys {
		value_item, _ := specs.GetObjectItem(key)
		if value_item == nil || value_item.IsNull() {
			continue
		}

		enabled, bool_err := value_item.GetBool()
		if bool_err != nil {
			return nil, sandbox.Deps.Std.Errorf("extension %q is not a bool", key)
		}

		conf.Extensions = append(conf.Extensions, Extension{
			Name:    key,
			Enabled: enabled,
		})
	}

	BindMethods(sandbox, conf)
	return conf, nil
}
