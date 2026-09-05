package resultconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

func New(deps *deps.Deps, content string) (*ResultConf, error) {

	if content == "" {
		return nil, deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsObject() {
		return nil, deps.Std.Errorf("result specs is not an object")
	}

	conf := NewEmpty(deps)

	if specs.HasKey("cli-output") {
		item, _ := specs.GetObjectItem("cli-output")
		if item != nil && item.IsString() {
			value, err := item.GetString()
			if err != nil {
				return nil, err
			}
			conf.CliOutput = value
		}
	}

	if specs.HasKey("exit-code") {
		item, _ := specs.GetObjectItem("exit-code")
		if item != nil && item.IsInt() {
			value, err := item.GetInt()
			if err != nil {
				return nil, err
			}
			conf.ExitCode = int(value)
		}
	}

	if specs.HasKey("tree") {
		if err := parseTree(deps, specs, conf); err != nil {
			return nil, err
		}
	}

	return conf, nil
}

// parseTree fills conf.Tree from the `tree` key: an array of
// {file, sha} objects. An entry missing either key is a malformed golden and
// is reported rather than silently compared as empty.
func parseTree(deps *deps.Deps, specs *serializibles.SerializibleObject, conf *ResultConf) error {
	tree, err := specs.GetObjectItem("tree")
	if err != nil {
		return err
	}
	if tree == nil || tree.IsNull() {
		return nil
	}
	if !tree.IsArray() {
		return deps.Std.Errorf("result tree is not an array")
	}

	size, err := tree.GetArraySize()
	if err != nil {
		return err
	}

	for index := 0; index < size; index++ {
		entry := tree.GetArrayItem(index)
		if entry == nil || !entry.IsObject() {
			return deps.Std.Errorf("result tree entry %d is not an object", index)
		}

		file, err := stringItem(deps, entry, "file", index)
		if err != nil {
			return err
		}
		sha, err := stringItem(deps, entry, "sha", index)
		if err != nil {
			return err
		}

		conf.AddTreeEntry(file, sha)
	}

	return nil
}

// stringItem reads one required string key of a tree entry, naming the entry
// when it is absent or of the wrong kind.
func stringItem(deps *deps.Deps, entry *serializibles.SerializibleObject, key string, index int) (string, error) {
	if !entry.HasKey(key) {
		return "", deps.Std.Errorf("result tree entry %d has no %s", index, key)
	}
	item, err := entry.GetObjectItem(key)
	if err != nil {
		return "", err
	}
	if item == nil || !item.IsString() {
		return "", deps.Std.Errorf("result tree entry %d has a non-string %s", index, key)
	}
	return item.GetString()
}
