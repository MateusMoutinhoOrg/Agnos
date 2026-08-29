package parsables

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

type PathReplacerEntry struct {
	Original    string
	Replacement string
}

type PathReplacerConf struct {
	Entries []PathReplacerEntry

	AddEntry func(original string, replacement string)
	Format   func(path string) string
	Render   func() string
}

func addPathReplacerConfMethods(sandbox *sandbox.SandBox, conf *PathReplacerConf) {

	conf.AddEntry = func(original string, replacement string) {
		conf.Entries = append(conf.Entries, PathReplacerEntry{
			Original:    original,
			Replacement: replacement,
		})
	}

	conf.Format = func(path string) string {
		result := path
		for _, entry := range conf.Entries {
			result = strings.ReplaceAll(result, entry.Original, entry.Replacement)
		}
		return result
	}

	conf.Render = func() string {
		obj := sandbox.Deps.SerializeLib.CreateObject()

		for _, entry := range conf.Entries {
			obj.AddItemToObject(entry.Original, entry.Replacement)
		}

		return sandbox.Deps.SerializeLib.SerializeToYaml(obj)
	}
}

func NewPathReplacerConf(sandbox *sandbox.SandBox, content string) (*PathReplacerConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Errorf("content cannot be empty, use NewPathReplacerConfEmpty instead")
	}

	specs, parse_error := sandbox.Deps.SerializeLib.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsObject() {
		return nil, sandbox.Deps.Errorf("paths config is not an object")
	}

	conf := &PathReplacerConf{
		Entries: make([]PathReplacerEntry, 0),
	}

	keys, keys_err := specs.GetKeys()
	if keys_err != nil {
		return nil, sandbox.Deps.Errorf("could not get paths keys")
	}

	for _, key := range keys {
		value_item, _ := specs.GetObjectItem(key)
		if value_item == nil || value_item.IsNull() {
			continue
		}

		value, str_err := value_item.GetString()
		if str_err != nil {
			return nil, sandbox.Deps.Errorf("path replacement for key %q is not a string", key)
		}

		conf.Entries = append(conf.Entries, PathReplacerEntry{
			Original:    key,
			Replacement: value,
		})
	}

	addPathReplacerConfMethods(sandbox, conf)
	return conf, nil
}

func NewPathReplacerConfEmpty(sandbox *sandbox.SandBox) *PathReplacerConf {
	conf := &PathReplacerConf{
		Entries: make([]PathReplacerEntry, 0),
	}
	addPathReplacerConfMethods(sandbox, conf)
	return conf
}
