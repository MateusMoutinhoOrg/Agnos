package opinionatedagnosserver

import (
	"sort"
	"strconv"
	"strings"

	serializabledeps "{{.Module}}/sandbox/deps/serializabledeps"
)

// validateForm converts the pairs of a form body into a document and checks it
// against schema_json, the canonical form of the route's declared
// form-schema — the flat subset of a json-schema, one property per form key.
// It is validateSchema for a body that arrives as `key=value` pairs rather
// than as json, and returns the same four values.
//
// Every value arrives as text, so each declared key is converted to what its
// property declares before the schema runs: an `integer` or a `number` is
// parsed, a `boolean` reads `true`/`1`/`on` and `false`/`0`/`off`, an `array`
// takes every occurrence of the key and any other type its first one. A value
// that will not convert stays text and is reported by the schema as of the
// wrong type. An empty value — an input left blank — counts as absent, so a
// `required` property catches it. A key the schema does not declare is kept as
// text, so `additionalProperties: false` still refuses it.
func validateForm(serializer serializabledeps.Contract, schema_json string, form map[string][]string) (*serializabledeps.SerializableObject, string, string, bool) {
	check := &validator{serializables: serializer}
	schema, err := check.serializables.ParseJson(schema_json)
	if err != nil {
		return nil, "", "the declared form-schema is not valid json", false
	}

	properties, _ := schema.GetObjectItem("properties")
	document := check.serializables.CreateObject()

	keys := make([]string, 0, len(form))
	for key := range form {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		present := formPresent(form[key])
		if len(present) == 0 {
			continue
		}

		var property *serializabledeps.SerializableObject
		if properties != nil && properties.IsObject() {
			property, _ = properties.GetObjectItem(key)
		}

		if property == nil {
			if len(present) == 1 {
				document.AddItemToObject(key, present[0])
				continue
			}
			document.AddItemToObject(key, check.formArray("", present))
			continue
		}

		declared := schemaString(property, "type")
		if declared == "array" {
			item_type := ""
			if items, _ := property.GetObjectItem("items"); items != nil {
				item_type = schemaString(items, "type")
			}
			document.AddItemToObject(key, check.formArray(item_type, present))
			continue
		}
		document.AddItemToObject(key, check.formValue(declared, present[0]))
	}

	field, message, ok := check.validateNode(schema, document, "")
	return document, field, message, ok
}

// formPresent drops the empty occurrences of one key: a blank input is sent as
// `key=`, and reads as if it had not been sent at all.
func formPresent(values []string) []string {
	present := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" {
			present = append(present, value)
		}
	}
	return present
}

// formArray builds the array one repeated key holds, each occurrence
// converted to the declared item type.
func (check *validator) formArray(declared string, values []string) *serializabledeps.SerializableObject {
	list := check.serializables.CreateArray()
	for _, value := range values {
		list.AddItemToArray(check.formValue(declared, value))
	}
	return list
}

// formValue converts one form value to the declared scalar type, returning
// the text itself when it will not convert — the schema then reports it.
func (check *validator) formValue(declared string, value string) any {
	switch declared {
	case "integer":
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	case "number":
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	case "boolean":
		switch strings.ToLower(value) {
		case "true", "1", "on":
			return true
		case "false", "0", "off":
			return false
		}
	}
	return value
}
