package routeio

import (
	serializables "{{.Module}}/sandbox/deps/serializables"
)

// The readers below are what the generated bindBody functions of every route
// convert a validated document with: one per Go type a json-schema property
// maps onto, each returning the zero value for a property that is absent or of
// another kind. The schema has already been enforced by ValidateSchema, so a
// reader never reports: what it cannot read was optional.

// ReadString returns the named property as text.
func ReadString(object *serializables.SerializibleObject, key string) string {
	item := childOf(object, key)
	if item == nil || !item.IsString() {
		return ""
	}
	value, err := item.GetString()
	if err != nil {
		return ""
	}
	return value
}

// ReadInt returns the named property as an int.
func ReadInt(object *serializables.SerializibleObject, key string) int {
	value, ok := numberValue(childOrNull(object, key))
	if !ok {
		return 0
	}
	return int(value)
}

// ReadFloat returns the named property as a float64.
func ReadFloat(object *serializables.SerializibleObject, key string) float64 {
	value, ok := numberValue(childOrNull(object, key))
	if !ok {
		return 0
	}
	return value
}

// ReadBool returns the named property as a bool.
func ReadBool(object *serializables.SerializibleObject, key string) bool {
	item := childOf(object, key)
	if item == nil || !item.IsBool() {
		return false
	}
	value, err := item.GetBool()
	if err != nil {
		return false
	}
	return value
}

// ReadObject returns the named property as a document to read further, or nil
// when it is absent.
func ReadObject(object *serializables.SerializibleObject, key string) *serializables.SerializibleObject {
	item := childOf(object, key)
	if item == nil || !item.IsObject() {
		return nil
	}
	return item
}

// ReadItems returns the named property's items, in order, for an array of any
// kind. An absent property yields no items.
func ReadItems(object *serializables.SerializibleObject, key string) []*serializables.SerializibleObject {
	item := childOf(object, key)
	if item == nil || !item.IsArray() {
		return nil
	}
	size, err := item.GetArraySize()
	if err != nil {
		return nil
	}
	items := make([]*serializables.SerializibleObject, 0, size)
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil {
			continue
		}
		items = append(items, entry)
	}
	return items
}

// ItemString reads one array item as text.
func ItemString(item *serializables.SerializibleObject) string {
	if item == nil || !item.IsString() {
		return ""
	}
	value, err := item.GetString()
	if err != nil {
		return ""
	}
	return value
}

// ItemInt reads one array item as an int.
func ItemInt(item *serializables.SerializibleObject) int {
	value, ok := numberValue(itemOrNull(item))
	if !ok {
		return 0
	}
	return int(value)
}

// ItemFloat reads one array item as a float64.
func ItemFloat(item *serializables.SerializibleObject) float64 {
	value, ok := numberValue(itemOrNull(item))
	if !ok {
		return 0
	}
	return value
}

// ItemBool reads one array item as a bool.
func ItemBool(item *serializables.SerializibleObject) bool {
	if item == nil || !item.IsBool() {
		return false
	}
	value, err := item.GetBool()
	if err != nil {
		return false
	}
	return value
}

// childOf returns the named property of object, or nil when either the object
// or the property is absent.
func childOf(object *serializables.SerializibleObject, key string) *serializables.SerializibleObject {
	if object == nil || !object.IsObject() {
		return nil
	}
	item, _ := object.GetObjectItem(key)
	if item == nil || item.IsNull() {
		return nil
	}
	return item
}

// childOrNull is childOf for the numeric readers, which take a node rather
// than an object and a key.
func childOrNull(object *serializables.SerializibleObject, key string) *serializables.SerializibleObject {
	item := childOf(object, key)
	if item == nil {
		return nullNode()
	}
	return item
}

// itemOrNull stands a missing array item in for a null one, so the numeric
// readers never take a nil.
func itemOrNull(item *serializables.SerializibleObject) *serializables.SerializibleObject {
	if item == nil {
		return nullNode()
	}
	return item
}

// nullNode is the stand-in a reader falls back to: a node that is of no kind,
// so every reader answers with its zero value.
func nullNode() *serializables.SerializibleObject {
	return &serializables.SerializibleObject{
		IsInt:   func() bool { return false },
		IsFloat: func() bool { return false },
	}
}
