package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

// Render serializes a RouteConf back to the route.yaml shape. Keys come out in
// one fixed order and comments are dropped, which is what makes a re-render
// idempotent — and why route.yaml is never edited by hand.
func Render(deps *deps.Deps, conf *RouteConf) string {
	obj := deps.Serializables.CreateObject()

	obj.AddItemToObject("method", conf.Method)
	obj.AddItemToObject("paths", segmentsArray(deps, conf.Paths))
	obj.AddItemToObject("category", conf.Category)
	obj.AddItemToObject("help", conf.Help)
	if conf.LongDescription != "" {
		obj.AddItemToObject("long-description", conf.LongDescription)
	}
	if len(conf.Examples) > 0 {
		obj.AddItemToObject("examples", stringArray(deps, conf.Examples))
	}
	if conf.Hidden {
		obj.AddItemToObject("hidden", true)
	}
	if len(conf.Headers) > 0 {
		obj.AddItemToObject("headers", fieldsArray(deps, conf.Headers))
	}
	if len(conf.Params) > 0 {
		obj.AddItemToObject("params", fieldsArray(deps, conf.Params))
	}
	if conf.Body.Type != BodyNone {
		obj.AddItemToObject("body", bodyObject(deps, conf.Body))
	}

	return deps.Serializables.SerializeToYaml(obj)
}

// SchemaJson renders the declared json-schema as canonical JSON: the same
// tree Render writes into route.yaml, serialized for the generated
// EntriesSchema constant routeio.ValidateSchema is run against.
func SchemaJson(deps *deps.Deps, conf *RouteConf) string {
	if !conf.Body.HasSchema || conf.Body.Schema == nil {
		return ""
	}
	return deps.Serializables.SerializeToJson(schemaObject(deps, conf.Body.Schema))
}

// segmentsArray renders `paths` as the ordered sequence it is: one entry per
// segment, a trigger carrying its `identifier` and a capture its field.
func segmentsArray(deps *deps.Deps, segments []Segment) *serializibles.SerializibleObject {
	arr := deps.Serializables.CreateArray()
	for _, segment := range segments {
		if segment.Field == nil {
			entry := deps.Serializables.CreateObject()
			entry.AddItemToObject("identifier", segment.Identifier)
			arr.AddItemToArray(entry)
			continue
		}
		arr.AddItemToArray(fieldObject(deps, *segment.Field))
	}
	return arr
}

// fieldsArray renders headers/params in the canonical ordered sequence shape.
func fieldsArray(deps *deps.Deps, fields []Field) *serializibles.SerializibleObject {
	arr := deps.Serializables.CreateArray()
	for _, field := range fields {
		arr.AddItemToArray(fieldObject(deps, field))
	}
	return arr
}

// fieldObject renders one field — a captured segment, a header or a query
// parameter — as its declaration entry.
func fieldObject(deps *deps.Deps, field Field) *serializibles.SerializibleObject {
	entry := deps.Serializables.CreateObject()
	entry.AddItemToObject("name", field.Key)
	if field.Description != "" {
		entry.AddItemToObject("description", field.Description)
	}
	if len(field.Examples) > 0 {
		entry.AddItemToObject("examples", stringArray(deps, field.Examples))
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
func bodyObject(deps *deps.Deps, body Body) *serializibles.SerializibleObject {
	entry := deps.Serializables.CreateObject()
	entry.AddItemToObject("type", body.Type)
	if body.Required {
		entry.AddItemToObject("required", true)
	}
	entry.AddItemToObject("max-bytes", int64(body.MaxBytes))
	if body.ContentType != "" {
		entry.AddItemToObject("content-type", body.ContentType)
	}
	if body.HasSchema && body.Schema != nil {
		entry.AddItemToObject("json-schema", schemaObject(deps, body.Schema))
	}
	return entry
}

// schemaObject renders one node of the declared json-schema, recursing through
// `properties` and `items`. Keys outside the subset were never parsed, so they
// are never written back.
func schemaObject(deps *deps.Deps, schema *Schema) *serializibles.SerializibleObject {
	entry := deps.Serializables.CreateObject()

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
		entry.AddItemToObject("required", stringArray(deps, schema.Required))
	}
	if schema.HasAdditionalProperties {
		entry.AddItemToObject("additionalProperties", schema.AdditionalProperties)
	}
	if len(schema.Enum) > 0 {
		entry.AddItemToObject("enum", stringArray(deps, schema.Enum))
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
		entry.AddItemToObject("items", schemaObject(deps, schema.Items))
	}
	if len(schema.Properties) > 0 {
		properties := deps.Serializables.CreateObject()
		for _, property := range schema.Properties {
			properties.AddItemToObject(property.Name, schemaObject(deps, property.Schema))
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

func stringArray(deps *deps.Deps, values []string) *serializibles.SerializibleObject {
	arr := deps.Serializables.CreateArray()
	for _, value := range values {
		arr.AddItemToArray(value)
	}
	return arr
}
