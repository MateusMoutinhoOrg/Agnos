package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/triggerconf"
)

// DefaultPriority is the rung `add-command` declares a command on, and
// DefaultMiddlewarePriority the one of a --middleware: low enough to run in
// front of every command, high enough to leave room below it.
const DefaultPriority = 100
const DefaultMiddlewarePriority = 10

// LastSegment is the End of an arg that runs through the last segment of the
// command line, whatever its length.
const LastSegment = -1

// ArgTypes is every type an arg may declare. Anything but a string is one
// segment long, and a segment that will not convert is a non-match.
var ArgTypes = []string{"string", "integer", "number", "uuid"}

// DefaultArgType is the type of an arg that declares none.
const DefaultArgType = "string"

// FlagTypes is every type a flag may declare.
var FlagTypes = []string{"string", "integer", "number", "boolean", "string-array", "integer-array"}

// DefaultFlagType is the type of a flag that declares none.
const DefaultFlagType = "string"

// TriggerTypes is every way a trigger compares a value.
var TriggerTypes = triggerconf.TriggerTypes

// legacyKeys are the keys an entries.yaml declared that command.yaml does not
// carry: `identifiers` on the command, `name`, `identifiers` and `array` on a
// field. A declaration still holding one is an old one.
var legacyKeys = []string{"identifiers"}
var legacyFieldKeys = []string{"name", "identifiers", "array", "examples"}

// New parses one command.yaml body into a CommandConf.
func New(sandbox *api.Sandbox, content string) (*CommandConf, error) {
	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := sandbox.Deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}
	if !specs.IsObject() {
		return nil, sandbox.Deps.Std.Errorf("command.yaml is not an object")
	}

	conf := NewEmpty(sandbox)
	conf.HasPriority = false
	conf.Priority = 0

	if item, _ := specs.GetObjectItem("priority"); item != nil && !item.IsNull() {
		conf.Priority = readInt(specs, "priority")
		conf.HasPriority = true
	}
	if item, _ := specs.GetObjectItem("segments"); item != nil && !item.IsNull() {
		conf.Segments = readInt(specs, "segments")
		conf.HasSegments = true
	}
	if item, _ := specs.GetObjectItem("strict"); item != nil && !item.IsNull() {
		conf.Strict = readBool(specs, "strict")
	}
	conf.Examples = readStringArray(specs, "examples")
	conf.Category = readString(specs, "category")
	conf.Help = readString(specs, "help")
	conf.LongDescription = readString(specs, "long-description")
	conf.Hidden = readBool(specs, "hidden")

	for _, key := range legacyKeys {
		if item, _ := specs.GetObjectItem(key); item != nil && !item.IsNull() {
			conf.Legacy = append(conf.Legacy, key)
		}
	}

	if item, _ := specs.GetObjectItem("args"); item != nil && item.IsArray() {
		args, legacy, err := readArgs(sandbox, item)
		if err != nil {
			return nil, err
		}
		conf.Args = args
		conf.Legacy = append(conf.Legacy, legacy...)
	}
	if item, _ := specs.GetObjectItem("flags"); item != nil && item.IsArray() {
		flags, legacy, err := readFlags(sandbox, item)
		if err != nil {
			return nil, err
		}
		conf.Flags = flags
		conf.Legacy = append(conf.Legacy, legacy...)
	}

	return conf, nil
}

// readArgs parses the `args` sequence, in declaration order.
func readArgs(sandbox *api.Sandbox, item *serializibles.SerializibleObject) ([]Arg, []string, error) {
	size, err := item.GetArraySize()
	if err != nil {
		return nil, nil, err
	}
	args := make([]Arg, 0, size)
	legacy := []string{}
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			continue
		}
		arg := Arg{
			Id:          readString(entry, "id"),
			Start:       readInt(entry, "start"),
			End:         readInt(entry, "end"),
			Type:        readString(entry, "type"),
			Required:    readBool(entry, "required"),
			Description: readString(entry, "description"),
			Trigger:     triggerconf.New(sandbox, entry),
		}
		if arg.Type == "" {
			arg.Type = DefaultArgType
		}
		if default_item, _ := entry.GetObjectItem("default"); default_item != nil && !default_item.IsNull() {
			arg.HasDefault = true
			arg.Default = anyToString(sandbox, default_item)
		}
		legacy = append(legacy, legacyOf(entry, "args")...)
		args = append(args, arg)
	}
	return args, legacy, nil
}

// readFlags parses the `flags` sequence, in declaration order.
func readFlags(sandbox *api.Sandbox, item *serializibles.SerializibleObject) ([]Flag, []string, error) {
	size, err := item.GetArraySize()
	if err != nil {
		return nil, nil, err
	}
	flags := make([]Flag, 0, size)
	legacy := []string{}
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			continue
		}
		flag := Flag{
			Id:          readString(entry, "id"),
			Type:        readString(entry, "type"),
			Required:    readBool(entry, "required"),
			Enum:        readStringArray(entry, "enum"),
			Pattern:     readString(entry, "pattern"),
			Description: readString(entry, "description"),
			Trigger:     triggerconf.New(sandbox, entry),
		}
		if keys_item, _ := entry.GetObjectItem("keys"); keys_item != nil && keys_item.IsArray() {
			flag.Keys = readStringArray(entry, "keys")
			flag.HasKeys = true
		} else {
			flag.Keys = []string{DefaultKey(sandbox, flag.Id)}
		}
		if flag.Type == "" {
			flag.Type = DefaultFlagType
		}
		if default_item, _ := entry.GetObjectItem("default"); default_item != nil && !default_item.IsNull() {
			flag.HasDefault = true
			flag.Default = anyToString(sandbox, default_item)
		}
		if min_item, _ := entry.GetObjectItem("min"); min_item != nil && !min_item.IsNull() {
			flag.Min, flag.HasMin = readNumber(min_item)
		}
		if max_item, _ := entry.GetObjectItem("max"); max_item != nil && !max_item.IsNull() {
			flag.Max, flag.HasMax = readNumber(max_item)
		}
		legacy = append(legacy, legacyOf(entry, "flags")...)
		flags = append(flags, flag)
	}
	return flags, legacy, nil
}

// legacyOf is every key of one arg or flag that only an entries.yaml declared.
func legacyOf(entry *serializibles.SerializibleObject, list string) []string {
	legacy := []string{}
	for _, key := range legacyFieldKeys {
		if item, _ := entry.GetObjectItem(key); item != nil && !item.IsNull() {
			legacy = append(legacy, list+"[]."+key)
		}
	}
	return legacy
}

// DefaultKey is the spelling a flag declaring no `keys` answers to: its id in
// kebab-case after "--" — Target is --target, OutFile is --out-file.
func DefaultKey(sandbox *api.Sandbox, id string) string {
	key := ""
	for index, letter := range id {
		if letter >= 'A' && letter <= 'Z' {
			if index > 0 {
				key += "-"
			}
			key += string(letter + ('a' - 'A'))
			continue
		}
		key += string(letter)
	}
	return "--" + key
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

func readInt(obj *serializibles.SerializibleObject, key string) int {
	item, _ := obj.GetObjectItem(key)
	if item == nil || item.IsNull() {
		return 0
	}
	value, ok := readNumber(item)
	if !ok {
		return 0
	}
	return int(value)
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

// anyToString renders a scalar yaml value as the text a default is spelled
// with in the generated Go literal.
func anyToString(sandbox *api.Sandbox, item *serializibles.SerializibleObject) string {
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
		return sandbox.Deps.Stringsdeps.FormatInt(value, 10)
	}
	if item.IsFloat() {
		value, _ := item.GetFloat()
		return sandbox.Deps.Stringsdeps.FormatFloat(value, 'g', -1, 64)
	}
	return ""
}
