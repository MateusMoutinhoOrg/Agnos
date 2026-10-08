package triggerconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializabledeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializabledeps"
)

// OneOf is the trigger type that holds on any of its `values`, which is how a
// command answers to more than one name.
const OneOf = "one-of"

// Regex is the trigger type whose value is a regular expression.
const Regex = "regex"

// TriggerTypes is every way a trigger compares a value, in the order docs
// list them.
var TriggerTypes = []string{"equal", "prefix", "text-prefix", "suffix", Regex, OneOf}

// New parses the `trigger` object of one entry; an entry with none comes back
// with Set false.
func New(sandbox *api.Sandbox, entry *serializabledeps.SerializableObject) Trigger {
	item, _ := entry.GetObjectItem("trigger")
	if item == nil || !item.IsObject() {
		return Trigger{}
	}
	return Trigger{
		Set:        true,
		Type:       readString(item, "type"),
		Value:      readString(item, "value"),
		Values:     readStringArray(item, "values"),
		Negate:     readBool(item, "negate"),
		IgnoreCase: readBool(item, "ignore-case"),
	}
}

func readString(obj *serializabledeps.SerializableObject, key string) string {
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

func readBool(obj *serializabledeps.SerializableObject, key string) bool {
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

func readStringArray(obj *serializabledeps.SerializableObject, key string) []string {
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
