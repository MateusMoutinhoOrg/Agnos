package routeio

import (
	"{{.Module}}/sandbox/api"
	serializables "{{.Module}}/sandbox/deps/serializables"
)

// Format patterns. The subset supports four `format` values, and each is
// enforced as the regular expression below through sandbox.Deps.Stringsdeps —
// deliberately pragmatic rather than a full grammar.
const (
	// FormatEmailPattern matches the shape of an email address.
	FormatEmailPattern = `^[^@\s]+@[^@\s.]+(\.[^@\s.]+)+$`
	// FormatUuidPattern matches a hyphenated uuid, either case.
	FormatUuidPattern = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	// FormatDateTimePattern matches an RFC 3339 timestamp.
	FormatDateTimePattern = `^\d{4}-\d{2}-\d{2}[Tt ]\d{2}:\d{2}:\d{2}(\.\d+)?([Zz]|[+-]\d{2}:\d{2})$`
	// FormatUriPattern matches an absolute uri.
	FormatUriPattern = `^[A-Za-z][A-Za-z0-9+.-]*:\S*$`
)

// ValidateSchema parses body as JSON and checks it against schema_json, the
// canonical form of the route's declared json-schema. It is a pure function
// over sandbox.Deps.Serializables: no route, no request and no response reach it.
//
// It returns the parsed document, the field path of the first violation, that
// violation's message, and whether the body passed. A body that passes comes
// back with an empty path and message; one that fails comes back with the
// document parsed as far as it got, so a caller may still report on it.
func ValidateSchema(sandbox *api.Sandbox, schema_json string, body []byte) (*serializables.SerializibleObject, string, string, bool) {
	parsed, err := sandbox.Deps.Serializables.ParseJson(string(body))
	if err != nil {
		return nil, "", "the request body is not valid json", false
	}

	if schema_json == "" {
		return parsed, "", "", true
	}

	schema, err := sandbox.Deps.Serializables.ParseJson(schema_json)
	if err != nil {
		return parsed, "", "the declared json-schema is not valid json", false
	}

	field, message, ok := validateNode(sandbox, schema, parsed, "")
	return parsed, field, message, ok
}

// validateNode checks one value against one schema node and recurses into its
// properties and items. It stops at the first violation: a caller answers 400
// with one message, so finding the rest would be work nobody reads.
func validateNode(sandbox *api.Sandbox, schema *serializables.SerializibleObject, value *serializables.SerializibleObject, path string) (string, string, bool) {
	if value == nil {
		return path, "is missing", false
	}

	if value.IsNull() && schemaBool(schema, "nullable") {
		return "", "", true
	}

	declared := schemaString(schema, "type")
	if declared != "" && !typeMatches(declared, value) {
		return path, "must be of type " + declared, false
	}

	if field, message, ok := validateConstAndEnum(sandbox, schema, value, path); !ok {
		return field, message, false
	}

	switch declared {
	case "object":
		return validateObject(sandbox, schema, value, path)
	case "array":
		return validateArray(sandbox, schema, value, path)
	case "string":
		return validateString(sandbox, schema, value, path)
	case "integer", "number":
		return validateNumber(sandbox, schema, value, path)
	}

	return "", "", true
}

// validateConstAndEnum enforces the two value-listing keywords, comparing the
// scalar renderings of the value and of each allowed entry.
func validateConstAndEnum(sandbox *api.Sandbox, schema *serializables.SerializibleObject, value *serializables.SerializibleObject, path string) (string, string, bool) {
	if item, _ := schema.GetObjectItem("const"); item != nil && !item.IsNull() {
		if scalarText(sandbox, item) != scalarText(sandbox, value) {
			return path, "must be " + scalarText(sandbox, item), false
		}
	}

	item, _ := schema.GetObjectItem("enum")
	if item == nil || !item.IsArray() {
		return "", "", true
	}

	size, err := item.GetArraySize()
	if err != nil {
		return "", "", true
	}

	text := scalarText(sandbox, value)
	allowed := make([]string, 0, size)
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil {
			continue
		}
		candidate := scalarText(sandbox, entry)
		if candidate == text {
			return "", "", true
		}
		allowed = append(allowed, candidate)
	}

	return path, "must be one of " + sandbox.Deps.Stringsdeps.Join(allowed, ", "), false
}

// validateObject enforces required, properties and additionalProperties.
func validateObject(sandbox *api.Sandbox, schema *serializables.SerializibleObject, value *serializables.SerializibleObject, path string) (string, string, bool) {
	for _, name := range schemaStringArray(schema, "required") {
		item, _ := value.GetObjectItem(name)
		if item == nil || item.IsNull() {
			return childPath(path, name), "is required", false
		}
	}

	properties, _ := schema.GetObjectItem("properties")

	if additional, _ := schema.GetObjectItem("additionalProperties"); additional != nil && additional.IsBool() {
		if allowed, _ := additional.GetBool(); !allowed {
			keys, err := value.GetKeys()
			if err == nil {
				sandbox.Deps.Sortdeps.Strings(keys)
				for _, key := range keys {
					if properties == nil || !properties.HasKey(key) {
						return childPath(path, key), "is not an allowed property", false
					}
				}
			}
		}
	}

	if properties == nil || !properties.IsObject() {
		return "", "", true
	}

	keys, err := properties.GetKeys()
	if err != nil {
		return "", "", true
	}
	sandbox.Deps.Sortdeps.Strings(keys)

	for _, key := range keys {
		item, _ := value.GetObjectItem(key)
		if item == nil || item.IsNull() {
			continue
		}
		property, _ := properties.GetObjectItem(key)
		if property == nil {
			continue
		}
		if field, message, ok := validateNode(sandbox, property, item, childPath(path, key)); !ok {
			return field, message, false
		}
	}

	return "", "", true
}

// validateArray enforces minItems, maxItems, uniqueItems and the items schema.
func validateArray(sandbox *api.Sandbox, schema *serializables.SerializibleObject, value *serializables.SerializibleObject, path string) (string, string, bool) {
	size, err := value.GetArraySize()
	if err != nil {
		return path, "is not a readable array", false
	}

	if bound, has := schemaNumber(schema, "minItems"); has && size < int(bound) {
		return path, "must hold at least " + sandbox.Deps.Stringsdeps.FormatInt(int64(bound), 10) + " item(s)", false
	}
	if bound, has := schemaNumber(schema, "maxItems"); has && size > int(bound) {
		return path, "must hold at most " + sandbox.Deps.Stringsdeps.FormatInt(int64(bound), 10) + " item(s)", false
	}

	if schemaBool(schema, "uniqueItems") {
		seen := map[string]bool{}
		for i := 0; i < size; i++ {
			entry := value.GetArrayItem(i)
			if entry == nil {
				continue
			}
			text := sandbox.Deps.Serializables.SerializeToJson(entry)
			if seen[text] {
				return path, "must hold unique items", false
			}
			seen[text] = true
		}
	}

	items, _ := schema.GetObjectItem("items")
	if items == nil || !items.IsObject() {
		return "", "", true
	}

	for i := 0; i < size; i++ {
		entry := value.GetArrayItem(i)
		if entry == nil {
			continue
		}
		index_path := path + "[" + sandbox.Deps.Stringsdeps.FormatInt(int64(i), 10) + "]"
		if field, message, ok := validateNode(sandbox, items, entry, index_path); !ok {
			return field, message, false
		}
	}

	return "", "", true
}

// validateString enforces minLength, maxLength, pattern and format. Lengths
// count runes, so a multi-byte character counts once.
func validateString(sandbox *api.Sandbox, schema *serializables.SerializibleObject, value *serializables.SerializibleObject, path string) (string, string, bool) {
	text, err := value.GetString()
	if err != nil {
		return path, "is not a readable string", false
	}
	length := len([]rune(text))

	if bound, has := schemaNumber(schema, "minLength"); has && length < int(bound) {
		return path, "must be at least " + sandbox.Deps.Stringsdeps.FormatInt(int64(bound), 10) + " character(s) long", false
	}
	if bound, has := schemaNumber(schema, "maxLength"); has && length > int(bound) {
		return path, "must be at most " + sandbox.Deps.Stringsdeps.FormatInt(int64(bound), 10) + " character(s) long", false
	}

	if pattern := schemaString(schema, "pattern"); pattern != "" {
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(pattern, text)
		if err != nil {
			return path, "is checked against an invalid pattern", false
		}
		if !matched {
			return path, "must match " + pattern, false
		}
	}

	format := schemaString(schema, "format")
	pattern, known := formatPattern(format)
	if !known {
		return "", "", true
	}

	matched, err := sandbox.Deps.Stringsdeps.MatchPattern(pattern, text)
	if err != nil || !matched {
		return path, "must be a valid " + format, false
	}

	return "", "", true
}

// validateNumber enforces the four numeric bounds.
func validateNumber(sandbox *api.Sandbox, schema *serializables.SerializibleObject, value *serializables.SerializibleObject, path string) (string, string, bool) {
	number, ok := numberValue(value)
	if !ok {
		return path, "is not a readable number", false
	}

	if bound, has := schemaNumber(schema, "minimum"); has && number < bound {
		return path, "must be >= " + numberText(sandbox, bound), false
	}
	if bound, has := schemaNumber(schema, "maximum"); has && number > bound {
		return path, "must be <= " + numberText(sandbox, bound), false
	}
	if bound, has := schemaNumber(schema, "exclusiveMinimum"); has && number <= bound {
		return path, "must be > " + numberText(sandbox, bound), false
	}
	if bound, has := schemaNumber(schema, "exclusiveMaximum"); has && number >= bound {
		return path, "must be < " + numberText(sandbox, bound), false
	}

	return "", "", true
}

// formatPattern is the regular expression one `format` value is enforced as,
// and whether the subset knows that value at all.
func formatPattern(format string) (string, bool) {
	switch format {
	case "email":
		return FormatEmailPattern, true
	case "uuid":
		return FormatUuidPattern, true
	case "date-time":
		return FormatDateTimePattern, true
	case "uri":
		return FormatUriPattern, true
	}
	return "", false
}

// typeMatches reports whether a parsed value is of the declared json type. An
// integer satisfies `number`, but a fractional value never satisfies
// `integer`.
func typeMatches(declared string, value *serializables.SerializibleObject) bool {
	switch declared {
	case "object":
		return value.IsObject()
	case "array":
		return value.IsArray()
	case "string":
		return value.IsString()
	case "boolean":
		return value.IsBool()
	case "integer":
		return value.IsInt()
	case "number":
		return value.IsInt() || value.IsFloat()
	case "null":
		return value.IsNull()
	}
	return true
}

// childPath is the field path of one property under path, the dotted spelling
// a caller reads in the "field" of the error body.
func childPath(path string, name string) string {
	if path == "" {
		return name
	}
	return path + "." + name
}

// scalarText renders one scalar as the text const and enum compare on.
func scalarText(sandbox *api.Sandbox, value *serializables.SerializibleObject) string {
	if value == nil {
		return ""
	}
	if value.IsString() {
		text, _ := value.GetString()
		return text
	}
	if value.IsBool() {
		if flag, _ := value.GetBool(); flag {
			return "true"
		}
		return "false"
	}
	if value.IsInt() {
		number, _ := value.GetInt()
		return sandbox.Deps.Stringsdeps.FormatInt(number, 10)
	}
	if value.IsFloat() {
		number, _ := value.GetFloat()
		return sandbox.Deps.Stringsdeps.FormatFloat(number, 'g', -1, 64)
	}
	return sandbox.Deps.Serializables.SerializeToJson(value)
}

// numberValue reads an int or a float node as one float64.
func numberValue(value *serializables.SerializibleObject) (float64, bool) {
	if value.IsInt() {
		number, err := value.GetInt()
		if err != nil {
			return 0, false
		}
		return float64(number), true
	}
	if value.IsFloat() {
		number, err := value.GetFloat()
		if err != nil {
			return 0, false
		}
		return number, true
	}
	return 0, false
}

// numberText spells a bound the way the declaration does: as an integer when
// it has no fractional part.
func numberText(sandbox *api.Sandbox, value float64) string {
	if value == float64(int64(value)) {
		return sandbox.Deps.Stringsdeps.FormatInt(int64(value), 10)
	}
	return sandbox.Deps.Stringsdeps.FormatFloat(value, 'g', -1, 64)
}

func schemaString(schema *serializables.SerializibleObject, key string) string {
	item, _ := schema.GetObjectItem(key)
	if item == nil || !item.IsString() {
		return ""
	}
	value, err := item.GetString()
	if err != nil {
		return ""
	}
	return value
}

func schemaBool(schema *serializables.SerializibleObject, key string) bool {
	item, _ := schema.GetObjectItem(key)
	if item == nil || !item.IsBool() {
		return false
	}
	value, err := item.GetBool()
	if err != nil {
		return false
	}
	return value
}

func schemaNumber(schema *serializables.SerializibleObject, key string) (float64, bool) {
	item, _ := schema.GetObjectItem(key)
	if item == nil {
		return 0, false
	}
	return numberValue(item)
}

func schemaStringArray(schema *serializables.SerializibleObject, key string) []string {
	item, _ := schema.GetObjectItem(key)
	if item == nil || !item.IsArray() {
		return nil
	}
	size, err := item.GetArraySize()
	if err != nil {
		return nil
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
