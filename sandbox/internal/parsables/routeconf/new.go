package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

// DefaultMethod is the http method a route that declares none answers to.
const DefaultMethod = "GET"

// DefaultMaxBytes is the body size limit a route that declares none carries,
// one mebibyte.
const DefaultMaxBytes = 1048576

// DefaultJsonContentType is the content type a `type: json` body that declares
// none requires.
const DefaultJsonContentType = "application/json"

// BodyNone is the body type of a route that takes no body — the default, and
// the one shape that generates no ReadBody at all.
const BodyNone = "none"

// New parses one route.yaml body into a RouteConf.
func New(sandbox *api.Sandbox, content string) (*RouteConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := sandbox.Deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsObject() {
		return nil, sandbox.Deps.Std.Errorf("route.yaml is not an object")
	}

	conf := &RouteConf{
		Paths:    []Segment{},
		Examples: []string{},
		Headers:  []Field{},
		Params:   []Field{},
		Body:     Body{Type: BodyNone, MaxBytes: DefaultMaxBytes},
	}

	conf.Method = normalizeMethod(sandbox, readString(specs, "method"))
	conf.Examples = readStringArray(specs, "examples")
	conf.Category = readString(specs, "category")
	conf.Help = readString(specs, "help")
	conf.LongDescription = readString(specs, "long-description")
	conf.Hidden = readBool(specs, "hidden")

	paths_item, _ := specs.GetObjectItem("paths")
	if paths_item != nil {
		segments, err := readSegments(sandbox, paths_item)
		if err != nil {
			return nil, err
		}
		conf.Paths = segments
	}

	headers_item, _ := specs.GetObjectItem("headers")
	if headers_item != nil {
		fields, err := readFieldCollection(sandbox, headers_item)
		if err != nil {
			return nil, err
		}
		conf.Headers = fields
	}

	params_item, _ := specs.GetObjectItem("params")
	if params_item != nil {
		fields, err := readFieldCollection(sandbox, params_item)
		if err != nil {
			return nil, err
		}
		conf.Params = fields
	}

	body_item, _ := specs.GetObjectItem("body")
	if body_item != nil && body_item.IsObject() {
		conf.Body = readBody(sandbox, body_item)
	}

	BindMethods(sandbox, conf)
	return conf, nil
}

// readSegments parses the `paths` sequence into the ordered segments the URL
// is matched against. Each entry carries either an `identifier` (a trigger) or
// a `name` (a capture); an entry carrying both, or neither, is refused here
// rather than generating a route nothing can match.
func readSegments(sandbox *api.Sandbox, item *serializibles.SerializibleObject) ([]Segment, error) {
	if item.IsNull() {
		return []Segment{}, nil
	}
	if !item.IsArray() {
		return nil, sandbox.Deps.Std.Errorf("`paths` must be a sequence of segments")
	}

	size, err := item.GetArraySize()
	if err != nil {
		return nil, err
	}

	segments := make([]Segment, 0, size)
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			return nil, sandbox.Deps.Std.Errorf("`paths` entry #%d is not an object", i)
		}

		identifier := readString(entry, "identifier")
		name := readString(entry, "name")

		if identifier != "" && name != "" {
			return nil, sandbox.Deps.Std.Errorf("`paths` entry #%d declares both `identifier` and `name`", i)
		}
		if identifier == "" && name == "" {
			return nil, sandbox.Deps.Std.Errorf("`paths` entry #%d needs an `identifier` or a `name`", i)
		}

		if identifier != "" {
			segments = append(segments, Segment{Identifier: identifier})
			continue
		}

		field := readFieldEntry(sandbox, entry)
		field.Key = name
		segments = append(segments, Segment{Field: &field})
	}

	return segments, nil
}

// readFieldCollection parses a headers/params declaration into an ordered
// slice of Field. Only the canonical YAML sequence is accepted: a header and a
// query key are external spellings, so each entry carries an explicit `name`.
func readFieldCollection(sandbox *api.Sandbox, item *serializibles.SerializibleObject) ([]Field, error) {
	if item.IsNull() {
		return []Field{}, nil
	}
	if !item.IsArray() {
		return nil, sandbox.Deps.Std.Errorf("`headers` and `params` must be sequences of fields")
	}

	size, err := item.GetArraySize()
	if err != nil {
		return nil, err
	}

	fields := make([]Field, 0, size)
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			continue
		}

		field := readFieldEntry(sandbox, entry)
		field.Key = readString(entry, "name")
		if field.Key == "" {
			return nil, sandbox.Deps.Std.Errorf("headers/params entry #%d needs a name", i)
		}
		fields = append(fields, field)
	}

	return fields, nil
}

// readFieldEntry parses the attributes shared by every field shape; the caller
// assigns Key.
func readFieldEntry(sandbox *api.Sandbox, item *serializibles.SerializibleObject) Field {
	field := Field{
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
		field.Default = anyToString(sandbox, default_item)
	}

	// `required` is meaningless for a boolean (absent means false) or for a
	// field that carries a default (the default covers its absence), so it is
	// dropped here and never reaches the generated dispatch.
	if field.Type == "boolean" || field.HasDefault {
		field.Required = false
	}

	return field
}

// readBody parses the `body` object, filling in the defaults a declaration
// leaves out: no body at all, one mebibyte, and application/json for json.
func readBody(sandbox *api.Sandbox, item *serializibles.SerializibleObject) Body {
	body := Body{
		Type:        normalizeBodyType(readString(item, "type")),
		Required:    readBool(item, "required"),
		MaxBytes:    DefaultMaxBytes,
		ContentType: readString(item, "content-type"),
	}

	if max_item, _ := item.GetObjectItem("max-bytes"); max_item != nil && !max_item.IsNull() {
		if value, ok := readNumber(max_item); ok {
			body.MaxBytes = int(value)
		}
	}

	if body.ContentType == "" && body.Type == "json" {
		body.ContentType = DefaultJsonContentType
	}

	if schema_item, _ := item.GetObjectItem("json-schema"); schema_item != nil && schema_item.IsObject() {
		body.Schema = readSchema(sandbox, schema_item)
		body.HasSchema = true
	}

	return body
}

// schemaKeys is the whole subset of JSON Schema a route may declare. Anything
// else ($ref, oneOf, allOf, anyOf, patternProperties) lands in Schema.Unknown,
// which build and verify report rather than silently drop.
var schemaKeys = []string{
	"type", "properties", "required", "additionalProperties", "items", "enum",
	"const", "minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum",
	"minLength", "maxLength", "pattern", "minItems", "maxItems", "uniqueItems",
	"format", "nullable",
}

// readSchema parses one node of the declared json-schema, recursing through
// `properties` and `items`. Property order is alphabetical, so the generated
// struct and the canonical schema JSON come out the same on every build.
func readSchema(sandbox *api.Sandbox, item *serializibles.SerializibleObject) *Schema {
	schema := &Schema{
		Type:        readString(item, "type"),
		Format:      readString(item, "format"),
		Pattern:     readString(item, "pattern"),
		Nullable:    readBool(item, "nullable"),
		UniqueItems: readBool(item, "uniqueItems"),
		Required:    readStringArray(item, "required"),
		Properties:  []SchemaProperty{},
	}

	schema.Unknown = unknownKeys(sandbox, item)

	if additional, _ := item.GetObjectItem("additionalProperties"); additional != nil && additional.IsBool() {
		schema.AdditionalProperties, _ = additional.GetBool()
		schema.HasAdditionalProperties = true
	}

	if const_item, _ := item.GetObjectItem("const"); const_item != nil && !const_item.IsNull() {
		schema.Const = anyToString(sandbox, const_item)
		schema.HasConst = true
	}

	if enum_item, _ := item.GetObjectItem("enum"); enum_item != nil && enum_item.IsArray() {
		size, err := enum_item.GetArraySize()
		if err == nil {
			for i := 0; i < size; i++ {
				entry := enum_item.GetArrayItem(i)
				if entry == nil {
					continue
				}
				schema.Enum = append(schema.Enum, anyToString(sandbox, entry))
			}
		}
	}

	schema.Minimum, schema.HasMinimum = readSchemaNumber(item, "minimum")
	schema.Maximum, schema.HasMaximum = readSchemaNumber(item, "maximum")
	schema.ExclusiveMinimum, schema.HasExclusiveMinimum = readSchemaNumber(item, "exclusiveMinimum")
	schema.ExclusiveMaximum, schema.HasExclusiveMaximum = readSchemaNumber(item, "exclusiveMaximum")

	schema.MinLength, schema.HasMinLength = readSchemaInt(item, "minLength")
	schema.MaxLength, schema.HasMaxLength = readSchemaInt(item, "maxLength")
	schema.MinItems, schema.HasMinItems = readSchemaInt(item, "minItems")
	schema.MaxItems, schema.HasMaxItems = readSchemaInt(item, "maxItems")

	if items_item, _ := item.GetObjectItem("items"); items_item != nil && items_item.IsObject() {
		schema.Items = readSchema(sandbox, items_item)
	}

	if properties, _ := item.GetObjectItem("properties"); properties != nil && properties.IsObject() {
		keys, err := properties.GetKeys()
		if err == nil {
			sandbox.Deps.Sortdeps.Strings(keys)
			for _, key := range keys {
				property, _ := properties.GetObjectItem(key)
				if property == nil || !property.IsObject() {
					continue
				}
				schema.Properties = append(schema.Properties, SchemaProperty{
					Name:   key,
					Schema: readSchema(sandbox, property),
				})
			}
		}
	}

	return schema
}

// unknownKeys lists the keys of one schema node that fall outside the subset,
// sorted so the message a build fails with is the same every time.
func unknownKeys(sandbox *api.Sandbox, item *serializibles.SerializibleObject) []string {
	keys, err := item.GetKeys()
	if err != nil {
		return nil
	}
	sandbox.Deps.Sortdeps.Strings(keys)

	var unknown []string
	for _, key := range keys {
		known := false
		for _, allowed := range schemaKeys {
			if allowed == key {
				known = true
				break
			}
		}
		if !known {
			unknown = append(unknown, key)
		}
	}
	return unknown
}

func readSchemaNumber(item *serializibles.SerializibleObject, key string) (float64, bool) {
	entry, _ := item.GetObjectItem(key)
	if entry == nil || entry.IsNull() {
		return 0, false
	}
	return readNumber(entry)
}

func readSchemaInt(item *serializibles.SerializibleObject, key string) (int, bool) {
	value, ok := readSchemaNumber(item, key)
	return int(value), ok
}

// normalizeMethod maps the method spellings accepted in route.yaml onto the
// upper-case set the dispatch compares against; "" is the default GET.
func normalizeMethod(sandbox *api.Sandbox, raw string) string {
	method := sandbox.Deps.Stringsdeps.ToUpper(sandbox.Deps.Stringsdeps.TrimSpace(raw))
	if method == "" {
		return DefaultMethod
	}
	return method
}

// normalizeBodyType maps the body-type spellings accepted in route.yaml onto
// the canonical set used everywhere else.
func normalizeBodyType(raw string) string {
	switch raw {
	case "raw", "bytes":
		return "raw"
	case "text", "string":
		return "text"
	case "json":
		return "json"
	default:
		return BodyNone
	}
}

// normalizeType maps the type spellings accepted in route.yaml onto the
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
