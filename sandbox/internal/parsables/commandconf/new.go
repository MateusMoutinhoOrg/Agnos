package commandconf

import (
	"strconv"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/serializebles"
)

// New parses one entries.yaml body into a CommandConf.
func New(deps *deps.Deps, content string) (*CommandConf, error) {

	if content == "" {
		return nil, deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := deps.Serializebles.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsObject() {
		return nil, deps.Std.Errorf("entries.yaml is not an object")
	}

	conf := &CommandConf{
		Identifiers: []string{},
		Examples:    []string{},
		Flags:       []Field{},
		Args:        []Field{},
	}

	conf.Identifiers = readStringArray(specs, "identifiers")
	conf.Examples = readStringArray(specs, "examples")
	conf.Category = readString(specs, "category")
	conf.Help = readString(specs, "help")
	conf.LongDescription = readString(specs, "long-description")
	conf.Hidden = readBool(specs, "hidden")

	flags_item, _ := specs.GetObjectItem("flags")
	if flags_item != nil && flags_item.IsObject() {
		fields, err := readFields(flags_item)
		if err != nil {
			return nil, err
		}
		conf.Flags = fields
	}

	args_item, _ := specs.GetObjectItem("args")
	if args_item != nil && args_item.IsObject() {
		fields, err := readFields(args_item)
		if err != nil {
			return nil, err
		}
		conf.Args = fields
	}

	BindMethods(deps, conf)
	return conf, nil
}

// readFields walks the ordered keys of a flags/args object and parses each
// entry into a Field.
func readFields(obj *serializibles.SerializibleObject) ([]Field, error) {
	keys, err := obj.GetKeys()
	if err != nil {
		return nil, err
	}

	fields := make([]Field, 0, len(keys))
	for _, key := range keys {
		item, _ := obj.GetObjectItem(key)
		if item == nil || !item.IsObject() {
			continue
		}

		field := Field{
			Key:         key,
			Identifiers: readStringArray(item, "identifiers"),
			Examples:    readStringArray(item, "examples"),
			Description: readString(item, "description"),
			Type:        normalizeType(readString(item, "type")),
			Required:    readBool(item, "required"),
			Array:       readBool(item, "array"),
			Min:         readInt(item, "min"),
			Max:         readInt(item, "max"),
		}

		if default_item, _ := item.GetObjectItem("default"); default_item != nil && !default_item.IsNull() {
			field.HasDefault = true
			field.Default = anyToString(default_item)
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// normalizeType maps the type spellings accepted in entries.yaml onto the
// canonical set used everywhere else.
func normalizeType(raw string) string {
	switch raw {
	case "bool", "boolean":
		return "boolean"
	case "int", "integer":
		return "int"
	case "float", "double", "number":
		return "float"
	default:
		return "string"
	}
}

func readString(obj *serializibles.SerializibleObject, key string) string {
	item, _ := obj.GetObjectItem(key)
	if item == nil || item.IsNull() {
		return ""
	}
	value, err := item.GetString()
	if err != nil {
		return ""
	}
	return value
}

func readBool(obj *serializibles.SerializibleObject, key string) bool {
	item, _ := obj.GetObjectItem(key)
	if item == nil || item.IsNull() {
		return false
	}
	value, err := item.GetBool()
	if err != nil {
		return false
	}
	return value
}

func readInt(obj *serializibles.SerializibleObject, key string) int {
	item, _ := obj.GetObjectItem(key)
	if item == nil || item.IsNull() {
		return 0
	}
	value, err := item.GetInt()
	if err != nil {
		return 0
	}
	return int(value)
}

func readStringArray(obj *serializibles.SerializibleObject, key string) []string {
	item, _ := obj.GetObjectItem(key)
	if item == nil || !item.IsArray() {
		return []string{}
	}
	size, err := item.GetArraySize()
	if err != nil {
		return []string{}
	}
	out := make([]string, 0, size)
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil {
			continue
		}
		value, err := entry.GetString()
		if err != nil {
			continue
		}
		out = append(out, value)
	}
	return out
}

// anyToString renders a scalar yaml value as the string that will be baked
// into the generated Go literal source.
func anyToString(item *serializibles.SerializibleObject) string {
	if item.IsString() {
		value, _ := item.GetString()
		return value
	}
	if item.IsBool() {
		if value, _ := item.GetBool(); value {
			return "true"
		}
		return "false"
	}
	if item.IsInt() {
		value, _ := item.GetInt()
		return strconv.FormatInt(value, 10)
	}
	if item.IsFloat() {
		value, _ := item.GetFloat()
		return strconv.FormatFloat(value, 'g', -1, 64)
	}
	return ""
}
