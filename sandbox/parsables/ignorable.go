package parsables

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

type IgnorableItens struct {
	Paths []string

	AddPath     func(path string)
	IsIgnorable func(path string) bool
	Render      func() string
}

func addIgnorableItensMethods(sandbox *api.SandBox, items *IgnorableItens) {

	items.AddPath = func(path string) {
		items.Paths = append(items.Paths, path)
	}

	items.IsIgnorable = func(path string) bool {
		for _, p := range items.Paths {
			if p == path {
				return true
			}
		}
		return false
	}

	items.Render = func() string {
		arr := sandbox.Deps.SerializeLib.CreateArray()

		for _, p := range items.Paths {
			arr.AddItemToArray(sandbox.Deps.SerializeLib.CreateString(p))
		}

		return sandbox.Deps.SerializeLib.SerializeToYaml(arr)
	}
}

func NewIgnorableItens(sandbox *api.SandBox, content string) (*IgnorableItens, error) {

	if content == "" {
		return nil, sandbox.Deps.Errorf("content cannot be empty, use NewIgnorableItensEmpty instead")
	}

	specs, parse_error := sandbox.Deps.SerializeLib.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsArray() {
		return nil, sandbox.Deps.Errorf("ignore config is not an array")
	}

	items := &IgnorableItens{
		Paths: make([]string, 0),
	}

	size, err := specs.GetArraySize()
	if err != nil {
		return nil, sandbox.Deps.Errorf("could not get ignore array size")
	}

	for i := 0; i < size; i++ {
		item := specs.GetArrayItem(i)
		if item == nil || item.IsNull() {
			continue
		}

		path, str_err := item.GetString()
		if str_err != nil {
			return nil, sandbox.Deps.Errorf("ignore item at index %d is not a string", i)
		}

		items.Paths = append(items.Paths, path)
	}

	addIgnorableItensMethods(sandbox, items)
	return items, nil
}

func NewIgnorableItensEmpty(sandbox *api.SandBox) *IgnorableItens {
	items := &IgnorableItens{
		Paths: make([]string, 0),
	}
	addIgnorableItensMethods(sandbox, items)
	return items
}
