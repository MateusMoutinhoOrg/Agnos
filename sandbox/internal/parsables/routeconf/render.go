package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

// Render serializes a RouteConf back to the route.yaml shape. Keys come out in
// one fixed order and comments are dropped, which is what makes a re-render
// idempotent — and why route.yaml is never edited by hand.
func Render(sandbox *api.Sandbox, conf *RouteConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()

	obj.AddItemToObject("method", conf.Method)
	obj.AddItemToObject("paths", segmentsArray(sandbox, conf.Paths))
	obj.AddItemToObject("category", conf.Category)
	obj.AddItemToObject("help", conf.Help)
	if conf.LongDescription != "" {
		obj.AddItemToObject("long-description", conf.LongDescription)
	}
	if len(conf.Examples) > 0 {
		obj.AddItemToObject("examples", stringArray(sandbox, conf.Examples))
	}
	if conf.Hidden {
		obj.AddItemToObject("hidden", true)
	}
	if len(conf.Headers) > 0 {
		obj.AddItemToObject("headers", fieldsArray(sandbox, conf.Headers))
	}
	if len(conf.Params) > 0 {
		obj.AddItemToObject("params", fieldsArray(sandbox, conf.Params))
	}
	if conf.Body.Type != BodyNone {
		obj.AddItemToObject("body", bodyObject(sandbox, conf.Body))
	}

	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}

// SchemaJson renders the declared json-schema as canonical JSON: the same
// tree Render writes into route.yaml, serialized for the generated
// EntriesSchema constant routeio.ValidateSchema is run against.
func SchemaJson(sandbox *api.Sandbox, conf *RouteConf) string {
	if !conf.Body.HasSchema || conf.Body.Schema == nil {
		return ""
	}
	return sandbox.Deps.Serializables.SerializeToJson(schemaObject(sandbox, conf.Body.Schema))
}

// segmentsArray renders `paths` as the ordered sequence it is: one entry per
// segment, a trigger carrying its `identifier` and a capture its field.
func segmentsArray(sandbox *api.Sandbox, segments []Segment) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, segment := range segments {
		if segment.Field == nil {
			entry := sandbox.Deps.Serializables.CreateObject()
			entry.AddItemToObject("identifier", segment.Identifier)
			arr.AddItemToArray(entry)
			continue
		}
		arr.AddItemToArray(fieldObject(sandbox, *segment.Field))
	}
	return arr
}

// fieldsArray renders headers/params in the canonical ordered sequence shape.
func fieldsArray(sandbox *api.Sandbox, fields []Field) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, field := range fields {
		arr.AddItemToArray(fieldObject(sandbox, field))
	}
	return arr
}

// fieldObject renders one field — a captured segment, a header or a query
// parameter — as its declaration entry.
func fieldObject(sandbox *api.Sandbox, field Field) *serializibles.SerializibleObject {
	entry := sandbox.Deps.Serializables.CreateObject()
	entry.AddItemToObject("name", field.Key)
	if field.Description != "" {
		entry.AddItemToObject("description", field.Description)
	}
	if len(field.Examples) > 0 {
		entry.AddItemToObject("examples", stringArray(sandbox, field.Examples))
	}
	entry.AddItemToObject("type", field.Type)
	if field.HasDefault {
		entry.AddItemToObject("default", field.Default)
	}
	if field.Required {
		entry.AddItemToObject("required", true)
	}
	if field.Array {
		entry.AddItemToObject("array", true)
	}
	if field.HasMin {
		entry.AddItemToObject("min", numberValue(field.Min))
	}
	if field.HasMax {
		entry.AddItemToObject("max", numberValue(field.Max))
	}
	return entry
}

// bodyObject renders the `body` declaration, the json-schema tree included.
func bodyObject(sandbox *api.Sandbox, body Body) *serializibles.SerializibleObject {
	entry := sandbox.Deps.Serializables.CreateObject()
	entry.AddItemToObject("type", body.Type)
	if body.Required {
		entry.AddItemToObject("required", true)
	}
	entry.AddItemToObject("max-bytes", int64(body.MaxBytes))
	if body.ContentType != "" {
		entry.AddItemToObject("content-type", body.ContentType)
	}
	if body.HasSchema && body.Schema != nil {
		entry.AddItemToObject("json-schema", schemaObject(sandbox, body.Schema))
	}
	return entry
}

// schemaObject renders one node of the declared json-schema, recursing through
// `properties` and `items`. Keys outside the subset were never parsed, so they
// are never written back.
func schemaObject(sandbox *api.Sandbox, schema *Schema) *serializibles.SerializibleObject {
	entry := sandbox.Deps.Serializables.CreateObject()

	if schema.Type != "" {
		entry.AddItemToObject("type", schema.Type)
	}
	if schema.Nullable {
		entry.AddItemToObject("nullable", true)
	}
	if schema.Format != "" {
		entry.AddItemToObject("format", schema.Format)
	}
	if schema.Pattern != "" {
		entry.AddItemToObject("pattern", schema.Pattern)
	}
	if len(schema.Required) > 0 {
		entry.AddItemToObject("required", stringArray(sandbox, schema.Required))
	}
	if schema.HasAdditionalProperties {
		entry.AddItemToObject("additionalProperties", schema.AdditionalProperties)
	}
	if len(schema.Enum) > 0 {
		entry.AddItemToObject("enum", stringArray(sandbox, schema.Enum))
	}
	if schema.HasConst {
		entry.AddItemToObject("const", schema.Const)
	}
	if schema.HasMinimum {
		entry.AddItemToObject("minimum", numberValue(schema.Minimum))
	}
	if schema.HasMaximum {
		entry.AddItemToObject("maximum", numberValue(schema.Maximum))
	}
	if schema.HasExclusiveMinimum {
		entry.AddItemToObject("exclusiveMinimum", numberValue(schema.ExclusiveMinimum))
	}
	if schema.HasExclusiveMaximum {
		entry.AddItemToObject("exclusiveMaximum", numberValue(schema.ExclusiveMaximum))
	}
	if schema.HasMinLength {
		entry.AddItemToObject("minLength", int64(schema.MinLength))
	}
	if schema.HasMaxLength {
		entry.AddItemToObject("maxLength", int64(schema.MaxLength))
	}
	if schema.HasMinItems {
		entry.AddItemToObject("minItems", int64(schema.MinItems))
	}
	if schema.HasMaxItems {
		entry.AddItemToObject("maxItems", int64(schema.MaxItems))
	}
	if schema.UniqueItems {
		entry.AddItemToObject("uniqueItems", true)
	}
	if schema.Items != nil {
		entry.AddItemToObject("items", schemaObject(sandbox, schema.Items))
	}
	if len(schema.Properties) > 0 {
		properties := sandbox.Deps.Serializables.CreateObject()
		for _, property := range schema.Properties {
			properties.AddItemToObject(property.Name, schemaObject(sandbox, property.Schema))
		}
		entry.AddItemToObject("properties", properties)
	}

	return entry
}

// numberValue writes a bound back as the scalar it reads as: an integer when
// the value has no fractional part, a float otherwise. Without it every bound
// would come back as "1" one build and "1.0" the next.
func numberValue(value float64) any {
	if value == float64(int64(value)) {
		return int64(value)
	}
	return value
}

func stringArray(sandbox *api.Sandbox, values []string) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, value := range values {
		arr.AddItemToArray(value)
	}
	return arr
}
