package ignorableconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func New(sandbox *api.Sandbox, content string) (*IgnorableConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := sandbox.Deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsArray() {
		return nil, sandbox.Deps.Std.Errorf("ignore config is not an array")
	}

	items := &IgnorableConf{
		Paths: make([]string, 0),
	}

	size, err := specs.GetArraySize()
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not get ignore array size")
	}

	for i := 0; i < size; i++ {
		item := specs.GetArrayItem(i)
		if item == nil || item.IsNull() {
			continue
		}

		path, str_err := item.GetString()
		if str_err != nil {
			return nil, sandbox.Deps.Std.Errorf("ignore item at index %d is not a string", i)
		}

		items.Paths = append(items.Paths, path)
	}

	BindMethods(sandbox, items)
	return items, nil
}
