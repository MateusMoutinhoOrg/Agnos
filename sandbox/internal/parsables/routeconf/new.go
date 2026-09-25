package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

// DefaultMethod is the http method `add-route` declares when it is given none.
const DefaultMethod = "GET"

// DefaultResponseType is the response-type `add-route` declares when it is
// given none.
const DefaultResponseType = "application/json"

// DefaultMaxBytes is the body size limit a route that declares none carries,
// one mebibyte.
const DefaultMaxBytes = 1048576

// DefaultJsonContentType is the content type a `type: json` body that declares
// none requires.
const DefaultJsonContentType = "application/json"

// BodyNone is the body type of a route that takes no body — the default, and
// the one shape that generates no ReadBody at all.
const BodyNone = "none"

// LastSegment is the End of a path slice that runs through the last segment
// of the request, whatever its length.
const LastSegment = -1

// TriggerTypes is every way a trigger compares a value, in the order docs
// list them.
var TriggerTypes = []string{"equal", "prefix", "suffix", "regex"}

// ParameterTypes is every type a parameter may declare.
var ParameterTypes = []string{"string", "number", "boolean", "datetime", "string-array"}

// ParameterFonts is every place a parameter may be read from.
var ParameterFonts = []string{"query", "header"}

// legacyKeys are the top-level keys of the declaration routeslist replaced.
var legacyKeys = []string{"method", "headers", "params"}

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
		Methods:    []string{},
		Paths:      []Path{},
		Parameters: []Parameter{},
		Examples:   []string{},
		Legacy:     []string{},
		Body:       Body{Type: BodyNone, MaxBytes: DefaultMaxBytes},
	}

	for _, method := range readStringArray(specs, "methods") {
		conf.Methods = append(conf.Methods, normalizeMethod(sandbox, method))
	}
	if priority_item, _ := specs.GetObjectItem("priority"); priority_item != nil && !priority_item.IsNull() {
		conf.Priority = readInt(specs, "priority")
		conf.HasPriority = true
	}
	conf.ResponseType = readString(specs, "response-type")
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

	paths_item, _ := specs.GetObjectItem("paths")
	if paths_item != nil {
		paths, err := readPaths(sandbox, paths_item)
		if err != nil {
			return nil, err
		}
		conf.Paths = paths
	}

	parameters_item, _ := specs.GetObjectItem("parameters")
	if parameters_item != nil {
		parameters, err := readParameters(sandbox, parameters_item)
		if err != nil {
			return nil, err
		}
		conf.Parameters = parameters
	}

	body_item, _ := specs.GetObjectItem("body")
	if body_item != nil && body_item.IsObject() {
		conf.Body = readBody(sandbox, body_item)
	}

	BindMethods(sandbox, conf)
	return conf, nil
}

// readPaths parses the `paths` sequence. Each entry needs an `id`; `start`
// defaults to 0 and `end` to the last segment, so an entry declaring neither
// reads the whole path.
func readPaths(sandbox *api.Sandbox, item *serializibles.SerializibleObject) ([]Path, error) {
	if item.IsNull() {
		return []Path{}, nil
	}
	if !item.IsArray() {
		return nil, sandbox.Deps.Std.Errorf("`paths` must be a sequence of path slices")
	}

	size, err := item.GetArraySize()
	if err != nil {
		return nil, err
	}

	paths := make([]Path, 0, size)
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			return nil, sandbox.Deps.Std.Errorf("`paths` entry #%d is not an object", i)
		}

		path := Path{
			Id:          readString(entry, "id"),
			Start:       0,
			End:         LastSegment,
			Trigger:     readTrigger(entry),
			Description: readString(entry, "description"),
		}
		if path.Id == "" {
			return nil, sandbox.Deps.Std.Errorf("`paths` entry #%d needs an `id`", i)
		}
		if start_item, _ := entry.GetObjectItem("start"); start_item != nil && !start_item.IsNull() {
			path.Start = readInt(entry, "start")
		}
		if end_item, _ := entry.GetObjectItem("end"); end_item != nil && !end_item.IsNull() {
			path.End = readInt(entry, "end")
		}
		paths = append(paths, path)
	}

	return paths, nil
}

// readParameters parses the `parameters` sequence. Each entry needs an `id`;
// `key` defaults to it, and `fonts` to the query string alone.
func readParameters(sandbox *api.Sandbox, item *serializibles.SerializibleObject) ([]Parameter, error) {
	if item.IsNull() {
		return []Parameter{}, nil
	}
	if !item.IsArray() {
		return nil, sandbox.Deps.Std.Errorf("`parameters` must be a sequence of parameters")
	}

	size, err := item.GetArraySize()
	if err != nil {
		return nil, err
	}

	parameters := make([]Parameter, 0, size)
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			return nil, sandbox.Deps.Std.Errorf("`parameters` entry #%d is not an object", i)
		}

		parameter := Parameter{
			Id:          readString(entry, "id"),
			Key:         readString(entry, "key"),
			Type:        normalizeType(readString(entry, "type")),
			Fonts:       readStringArray(entry, "fonts"),
			Required:    readBool(entry, "required"),
			Trigger:     readTrigger(entry),
			Description: readString(entry, "description"),
			Examples:    readStringArray(entry, "examples"),
		}
		if parameter.Id == "" {
			return nil, sandbox.Deps.Std.Errorf("`parameters` entry #%d needs an `id`", i)
		}
		if parameter.Key == "" {
			parameter.Key = parameter.Id
		}
		if default_item, _ := entry.GetObjectItem("default"); default_item != nil && !default_item.IsNull() {
			parameter.HasDefault = true
			parameter.Default = anyToString(sandbox, default_item)
		}
		parameters = append(parameters, parameter)
	}

	return parameters, nil
}

// readTrigger parses the `trigger` object of a path or a parameter; an entry
// with none comes back with Exists false.
func readTrigger(entry *serializibles.SerializibleObject) Trigger {
	item, _ := entry.GetObjectItem("trigger")
	if item == nil || !item.IsObject() {
		return Trigger{}
	}
	return Trigger{
		Exists: true,
		Type:   readString(item, "type"),
		Value:  readString(item, "value"),
	}
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
// upper-case set the dispatch compares against.
func normalizeMethod(sandbox *api.Sandbox, raw string) string {
	return sandbox.Deps.Stringsdeps.ToUpper(sandbox.Deps.Stringsdeps.TrimSpace(raw))
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
// canonical set used everywhere else; an unknown spelling is kept, so verify
// can name it.
func normalizeType(raw string) string {
	switch raw {
	case "", "string":
		return "string"
	case "bool", "boolean":
		return "boolean"
	case "int", "integer", "float", "double", "number":
		return "number"
	default:
		return raw
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

// readInt reads a whole-number key as an int, answering 0 for an absent, null
// or non-numeric one — the same shape readString and readBool have.
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
