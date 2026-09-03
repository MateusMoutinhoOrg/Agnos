package depsversionconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func New(deps *deps.Deps, content string) (*DepsVersionConf, error) {

	if content == "" {
		return nil, deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsObject() {
		return nil, deps.Std.Errorf("depsversion specs is not an object")
	}

	keys, err := specs.GetKeys()
	if err != nil {
		return nil, err
	}

	conf := &DepsVersionConf{
		Deps: map[string]string{},
	}

	for _, key := range keys {
		item, _ := specs.GetObjectItem(key)
		if item == nil || item.IsNull() {
			continue
		}
		value, err := item.GetString()
		if err != nil {
			return nil, deps.Std.Errorf("depsversion entry %q is not a string", key)
		}
		conf.Deps[key] = value
	}

	BindMethods(deps, conf)
	return conf, nil
}
