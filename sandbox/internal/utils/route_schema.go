package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
)

// The json-schema helpers below are to a route's `body` what the field helpers
// of route_conf.go are to its `headers` and `params`: the shape both
// `add-body-field` and `remove-body-field` read a declaration by, so a
// property is found, inserted and dropped the same way on either side.

// RouteBodyTypes is every value a route's `body.type` may spell.
var RouteBodyTypes = []string{"none", "raw", "text", "json"}

// RouteBodyType normalizes a body type, taking the aliases the parser takes
// ("bytes" for raw, "string" for text) and refusing anything else.
func RouteBodyType(sandbox *api.Sandbox, raw string) (string, error) {
	switch sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw)) {
	case routeconf.BodyNone:
		return routeconf.BodyNone, nil
	case "raw", "bytes":
		return "raw", nil
	case "text", "string":
		return "text", nil
	case "json":
		return "json", nil
	}
	return "", sandbox.Deps.Std.Errorf("unknown body type %q (use %s)", raw, sandbox.Deps.Stringsdeps.Join(RouteBodyTypes, ", "))
}

// RouteSchemaFormats is every `format` a string property may declare.
var RouteSchemaFormats = []string{"email", "uuid", "date-time", "uri"}

// RouteSchemaFormat normalizes a --format, listing the accepted values when it
// is not one of them.
func RouteSchemaFormat(sandbox *api.Sandbox, raw string) (string, error) {
	format := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw))
	if format == "" {
		return "", nil
	}
	for _, known := range RouteSchemaFormats {
		if known == format {
			return format, nil
		}
	}
	return "", sandbox.Deps.Std.Errorf("unknown format %q (use %s)", raw, sandbox.Deps.Stringsdeps.Join(RouteSchemaFormats, ", "))
}

// RouteSchemaKind normalizes the --type of a body property: the four field
// types a request line carries, plus "object" — which a body has and a header
// or a query key cannot.
func RouteSchemaKind(sandbox *api.Sandbox, raw string) (string, error) {
	if sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw)) == "object" {
		return "object", nil
	}
	kind, ok := FieldType(sandbox, raw)
	if !ok {
		return "", sandbox.Deps.Std.Errorf("unknown type %q (use string, boolean, int, float or object)", raw)
	}
	return kind, nil
}

// RouteSchemaType maps a declared property type onto its json-schema spelling.
func RouteSchemaType(kind string) string {
	switch kind {
	case "int":
		return "integer"
	case "float":
		return "number"
	case "boolean":
		return "boolean"
	case "object":
		return "object"
	default:
		return "string"
	}
}

// SplitSchemaPath breaks a dotted property path ("address.city") into the
// segments it names.
func SplitSchemaPath(sandbox *api.Sandbox, name string) []string {
	return sandbox.Deps.Stringsdeps.Split(RouteFieldName(sandbox, name), ".")
}

// SchemaPropertyOf returns the named property of an object schema, or nil.
func SchemaPropertyOf(schema *routeconf.Schema, name string) *routeconf.Schema {
	if schema == nil {
		return nil
	}
	for _, property := range schema.Properties {
		if property.Name == name {
			return property.Schema
		}
	}
	return nil
}

// SchemaObjectOf is the object a declaration carries, reached through the
// array wrapping it when there is one: a property is walked into the same way
// whether it was declared as an object or as an array of them.
func SchemaObjectOf(schema *routeconf.Schema) *routeconf.Schema {
	if schema != nil && schema.Type == "array" {
		return schema.Items
	}
	return schema
}

// InsertSchemaProperty adds one property in alphabetical position, which is
// the order the declaration is read back in — so the generated struct and the
// rendered route.yaml are the same on every build.
func InsertSchemaProperty(parent *routeconf.Schema, name string, schema *routeconf.Schema) {
	entry := routeconf.SchemaProperty{Name: name, Schema: schema}
	for i, property := range parent.Properties {
		if name < property.Name {
			properties := make([]routeconf.SchemaProperty, 0, len(parent.Properties)+1)
			properties = append(properties, parent.Properties[:i]...)
			properties = append(properties, entry)
			properties = append(properties, parent.Properties[i:]...)
			parent.Properties = properties
			return
		}
	}
	parent.Properties = append(parent.Properties, entry)
}

// DropSchemaProperty removes the named property of an object schema and
// unlists it from that object's `required`, reporting whether it was declared
// at all.
func DropSchemaProperty(parent *routeconf.Schema, name string) bool {
	index := -1
	for i, property := range parent.Properties {
		if property.Name == name {
			index = i
			break
		}
	}
	if index < 0 {
		return false
	}

	properties := make([]routeconf.SchemaProperty, 0, len(parent.Properties)-1)
	properties = append(properties, parent.Properties[:index]...)
	properties = append(properties, parent.Properties[index+1:]...)
	parent.Properties = properties

	required := make([]string, 0, len(parent.Required))
	for _, entry := range parent.Required {
		if entry != name {
			required = append(required, entry)
		}
	}
	parent.Required = required

	return true
}
