package pathreplacerconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func New(sandbox *api.Sandbox, content string) (*PathReplacerConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := sandbox.Deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsObject() {
		return nil, sandbox.Deps.Std.Errorf("paths config is not an object")
	}

	conf := &PathReplacerConf{
		Entries: make([]PathReplacerEntry, 0),
	}

	keys, keys_err := specs.GetKeys()
	if keys_err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not get paths keys")
	}

	for _, key := range keys {
		value_item, _ := specs.GetObjectItem(key)
		if value_item == nil || value_item.IsNull() {
			continue
		}

		value, str_err := value_item.GetString()
		if str_err != nil {
			return nil, sandbox.Deps.Std.Errorf("path replacement for key %q is not a string", key)
		}

		conf.Entries = append(conf.Entries, PathReplacerEntry{
			Original:    key,
			Replacement: value,
		})
	}

	BindMethods(sandbox, conf)
	return conf, nil
}
