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

	obj.AddItemToObject("methods", stringArray(sandbox, conf.Methods))
	obj.AddItemToObject("priority", int64(conf.Priority))
	obj.AddItemToObject("response-type", conf.ResponseType)
	if conf.HasSegments {
		obj.AddItemToObject("segments", int64(conf.Segments))
	}
	if conf.Phase != "" && conf.Phase != PhaseBefore {
		obj.AddItemToObject("phase", conf.Phase)
	}
	obj.AddItemToObject("paths", pathsArray(sandbox, conf.Paths))
	if len(conf.Parameters) > 0 {
		obj.AddItemToObject("parameters", parametersArray(sandbox, conf.Parameters))
	}
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
	if conf.Body.Type != BodyNone {
		obj.AddItemToObject("body", bodyObject(sandbox, conf.Body))
	}

	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}

// SchemaJson renders the declared json-schema as canonical JSON: the same
// tree Render writes into route.yaml, serialized for the generated
// BodySchema constant routeio.ValidateSchema is run against.
func SchemaJson(sandbox *api.Sandbox, conf *RouteConf) string {
	if !conf.Body.HasSchema || conf.Body.Schema == nil {
		return ""
	}
	return sandbox.Deps.Serializables.SerializeToJson(schemaObject(sandbox, conf.Body.Schema))
}

// pathsArray renders `paths` as the ordered sequence it is: one entry per
// path slice, its trigger included when it declares one.
func pathsArray(sandbox *api.Sandbox, paths []Path) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, path := range paths {
		entry := sandbox.Deps.Serializables.CreateObject()
		entry.AddItemToObject("id", path.Id)
		if path.Description != "" {
			entry.AddItemToObject("description", path.Description)
		}
		entry.AddItemToObject("start", int64(path.Start))
		entry.AddItemToObject("end", int64(path.End))
		if path.Type != "" && path.Type != DefaultPathType {
			entry.AddItemToObject("type", path.Type)
		}
		if path.Trigger.Exists {
			entry.AddItemToObject("trigger", triggerObject(sandbox, path.Trigger))
		}
		arr.AddItemToArray(entry)
	}
	return arr
}

// parametersArray renders `parameters` in declaration order. `key` is written
// only when it differs from `id`, the one spelling it defaults to.
func parametersArray(sandbox *api.Sandbox, parameters []Parameter) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, parameter := range parameters {
		entry := sandbox.Deps.Serializables.CreateObject()
		entry.AddItemToObject("id", parameter.Id)
		if parameter.Key != "" && parameter.Key != parameter.Id {
			entry.AddItemToObject("key", parameter.Key)
		}
		if parameter.Description != "" {
			entry.AddItemToObject("description", parameter.Description)
		}
		if len(parameter.Examples) > 0 {
			entry.AddItemToObject("examples", stringArray(sandbox, parameter.Examples))
		}
		entry.AddItemToObject("type", parameter.Type)
		entry.AddItemToObject("fonts", stringArray(sandbox, parameter.Fonts))
		if parameter.Required {
			entry.AddItemToObject("required", true)
		}
		if parameter.HasDefault {
			entry.AddItemToObject("default", parameter.Default)
		}
		if parameter.Trigger.Exists {
			entry.AddItemToObject("trigger", triggerObject(sandbox, parameter.Trigger))
		}
		arr.AddItemToArray(entry)
	}
	return arr
}

// triggerObject renders one `trigger`.
func triggerObject(sandbox *api.Sandbox, trigger Trigger) *serializibles.SerializibleObject {
	entry := sandbox.Deps.Serializables.CreateObject()
	entry.AddItemToObject("type", trigger.Type)
	entry.AddItemToObject("value", trigger.Value)
	if trigger.Negate {
		entry.AddItemToObject("negate", true)
	}
	if trigger.IgnoreCase {
		entry.AddItemToObject("ignore-case", true)
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
