package commandconf

import (
	"sort"
	"strconv"
	"strings"

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
	if flags_item != nil {
		fields, err := readFieldCollection(deps, flags_item)
		if err != nil {
			return nil, err
		}
		conf.Flags = fields
	}

	args_item, _ := specs.GetObjectItem("args")
	if args_item != nil {
		fields, err := readFieldCollection(deps, args_item)
		if err != nil {
			return nil, err
		}
		conf.Args = fields
	}

	BindMethods(deps, conf)
	return conf, nil
}

// readFieldCollection parses a flags/args declaration into an ordered slice of
// Field. Two shapes are accepted:
//
//   - a YAML sequence (canonical) — each entry an object carrying an explicit
//     `name` (or, for flags, deriving one from the first long `--identifier`).
//     Sequences are ordered, so positional args bind by written position.
//   - a YAML mapping (legacy) — the key names the field. Mappings are
//     unordered, so the keys are sorted for a deterministic result; declare
//     positional args as a sequence when order matters.
func readFieldCollection(deps *deps.Deps, item *serializibles.SerializibleObject) ([]Field, error) {
	if item.IsArray() {
		return readFieldsFromArray(deps, item)
	}
	if item.IsObject() {
		return readFieldsFromObject(item)
	}
	return []Field{}, nil
}

func readFieldsFromArray(deps *deps.Deps, arr *serializibles.SerializibleObject) ([]Field, error) {
	size, err := arr.GetArraySize()
	if err != nil {
		return nil, err
	}

	fields := make([]Field, 0, size)
	for i := 0; i < size; i++ {
		entry := arr.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			continue
		}

		field := readFieldEntry(entry)
		field.Key = fieldKey(entry, field.Identifiers)
		if field.Key == "" {
			return nil, deps.Std.Errorf("flags/args entry #%d needs a name (or a -- identifier)", i)
		}
		fields = append(fields, field)
	}

	return fields, nil
}

func readFieldsFromObject(obj *serializibles.SerializibleObject) ([]Field, error) {
	keys, err := obj.GetKeys()
	if err != nil {
		return nil, err
	}
	sort.Strings(keys)

	fields := make([]Field, 0, len(keys))
	for _, key := range keys {
		item, _ := obj.GetObjectItem(key)
		if item == nil || !item.IsObject() {
			continue
		}
		field := readFieldEntry(item)
		field.Key = key
		fields = append(fields, field)
	}

	return fields, nil
}

// readFieldEntry parses the attributes common to both shapes; the caller
// assigns Key.
func readFieldEntry(item *serializibles.SerializibleObject) Field {
	field := Field{
		Identifiers: readStringArray(item, "identifiers"),
		Examples:    readStringArray(item, "examples"),
		Description: readString(item, "description"),
		Type:        normalizeType(readString(item, "type")),
		Required:    readBool(item, "required"),
		Array:       readBool(item, "array"),
	}

	if min_item, _ := item.GetObjectItem("min"); min_item != nil && !min_item.IsNull() {
		field.Min, field.HasMin = readNumber(min_item)
	}
	if max_item, _ := item.GetObjectItem("max"); max_item != nil && !max_item.IsNull() {
		field.Max, field.HasMax = readNumber(max_item)
	}

	if default_item, _ := item.GetObjectItem("default"); default_item != nil && !default_item.IsNull() {
		field.HasDefault = true
		field.Default = anyToString(default_item)
	}

	// `required` is meaningless for a boolean (absent means false) or for a
	// field that carries a default (the default covers its absence), so it is
	// dropped here and never reaches the generated dispatch or help.
	if field.Type == "boolean" || field.HasDefault {
		field.Required = false
	}

	return field
}

// fieldKey resolves the generated struct field name for a sequence entry:
// an explicit `name`, else the first long `--identifier` with its dashes
// stripped, else the first identifier.
func fieldKey(entry *serializibles.SerializibleObject, identifiers []string) string {
	if name := readString(entry, "name"); name != "" {
		return name
	}
	for _, id := range identifiers {
		if strings.HasPrefix(id, "--") {
			return strings.TrimLeft(id, "-")
		}
	}
	for _, id := range identifiers {
		return strings.TrimLeft(id, "-")
	}
	return ""
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

// readNumber reads a scalar yaml int or float as a float64, reporting whether
// a usable numeric value was found.
func readNumber(item *serializibles.SerializibleObject) (float64, bool) {
	if item.IsInt() {
		value, err := item.GetInt()
		if err != nil {
			return 0, false
		}
		return float64(value), true
	}
	if item.IsFloat() {
		value, err := item.GetFloat()
		if err != nil {
			return 0, false
		}
		return value, true
	}
	return 0, false
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
