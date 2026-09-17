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

// WalkSchemaTo descends the dotted path to the object the leaf is declared in,
// creating the objects it passes through. An array of objects is entered
// through its `items`, so "tags.name" reads the same whether tags is an object
// or a list of them.
func WalkSchemaTo(sandbox *api.Sandbox, root *routeconf.Schema, segments []string) (*routeconf.Schema, error) {
	parent := root
	for _, segment := range segments {
		property := SchemaPropertyOf(parent, segment)
		if property == nil {
			child := &routeconf.Schema{Type: "object"}
			InsertSchemaProperty(parent, segment, child)
			parent = child
			continue
		}

		child := SchemaObjectOf(property)
		if child == nil || child.Type != "object" {
			return nil, sandbox.Deps.Std.Errorf("body property %q is not an object, so nothing can be declared under it", segment)
		}
		parent = child
	}
	return parent, nil
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

// RouteBodyPropertySchema turns the raw values typed on the command line into
// one json-schema node. --array wraps the whole of it, so every keyword but
// the three list ones describes the element and not the list.
//
// It is the one reading of those flags: `add-body-field` builds a property
// with it, and `set-body-field` rebuilds one with it from the keywords already
// declared plus the ones being changed — so a property that was edited and one
// that was declared outright come out the same.
func RouteBodyPropertySchema(sandbox *api.Sandbox, props api.RouteBodyFieldProps) (*routeconf.Schema, error) {
	kind, err := RouteSchemaKind(sandbox, props.Type)
	if err != nil {
		return nil, err
	}
	format, err := RouteSchemaFormat(sandbox, props.Format)
	if err != nil {
		return nil, err
	}

	leaf := &routeconf.Schema{
		Type:     RouteSchemaType(kind),
		Nullable: props.Nullable,
		Format:   format,
		Pattern:  sandbox.Deps.Stringsdeps.TrimSpace(props.Pattern),
		Enum:     props.Enum,
	}
	if (leaf.Format != "" || leaf.Pattern != "") && kind != "string" {
		return nil, sandbox.Deps.Std.Errorf("--format and --pattern only apply to a string property")
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Const); value != "" {
		leaf.Const, leaf.HasConst = value, true
	}

	if err := applySchemaBounds(sandbox, leaf, kind, props); err != nil {
		return nil, err
	}
	if err := applySchemaObjectKeys(sandbox, leaf, kind, props); err != nil {
		return nil, err
	}

	return applySchemaListKeys(sandbox, leaf, props)
}

// applySchemaBounds spells --min/--max as the keywords the property's type
// uses: minimum/maximum for a number, minLength/maxLength for text. The two
// exclusive bounds are a number's alone.
func applySchemaBounds(sandbox *api.Sandbox, schema *routeconf.Schema, kind string, props api.RouteBodyFieldProps) error {
	numeric := kind == "int" || kind == "float"

	if props.ExclusiveMin != "" || props.ExclusiveMax != "" {
		if !numeric {
			return sandbox.Deps.Std.Errorf("--exclusive-min and --exclusive-max only apply to an int or float property")
		}
		if props.ExclusiveMin != "" {
			value, err := RouteParseBound(sandbox, kind, "exclusive-min", props.ExclusiveMin)
			if err != nil {
				return err
			}
			schema.ExclusiveMinimum, schema.HasExclusiveMinimum = value, true
		}
		if props.ExclusiveMax != "" {
			value, err := RouteParseBound(sandbox, kind, "exclusive-max", props.ExclusiveMax)
			if err != nil {
				return err
			}
			schema.ExclusiveMaximum, schema.HasExclusiveMaximum = value, true
		}
	}

	if props.Min == "" && props.Max == "" {
		return nil
	}
	if !numeric && kind != "string" {
		return sandbox.Deps.Std.Errorf("min/max do not apply to a %s property", kind)
	}

	if kind == "string" {
		if props.Min != "" {
			value, err := sandbox.Deps.Stringsdeps.Atoi(props.Min)
			if err != nil {
				return sandbox.Deps.Std.Errorf("min must be an int, got %q", props.Min)
			}
			schema.MinLength, schema.HasMinLength = value, true
		}
		if props.Max != "" {
			value, err := sandbox.Deps.Stringsdeps.Atoi(props.Max)
			if err != nil {
				return sandbox.Deps.Std.Errorf("max must be an int, got %q", props.Max)
			}
			schema.MaxLength, schema.HasMaxLength = value, true
		}
		return nil
	}

	if props.Min != "" {
		value, err := RouteParseBound(sandbox, kind, "min", props.Min)
		if err != nil {
			return err
		}
		schema.Minimum, schema.HasMinimum = value, true
	}
	if props.Max != "" {
		value, err := RouteParseBound(sandbox, kind, "max", props.Max)
		if err != nil {
			return err
		}
		schema.Maximum, schema.HasMaximum = value, true
	}
	return nil
}

// applySchemaObjectKeys settles whether an object property accepts keys it
// does not declare. It is the one switch that has to be written even when it
// is false, so the schema tells "undeclared keys are refused" from "nothing
// was said".
func applySchemaObjectKeys(sandbox *api.Sandbox, schema *routeconf.Schema, kind string, props api.RouteBodyFieldProps) error {
	if !props.AdditionalProperties && !props.NoAdditionalProperties {
		return nil
	}
	if props.AdditionalProperties && props.NoAdditionalProperties {
		return sandbox.Deps.Std.Errorf("--additional-properties and --no-additional-properties are mutually exclusive")
	}
	if kind != "object" {
		return sandbox.Deps.Std.Errorf("--additional-properties only applies to an object property")
	}

	schema.AdditionalProperties = props.AdditionalProperties
	schema.HasAdditionalProperties = true
	return nil
}

// applySchemaListKeys wraps the element schema in the array --array asks for,
// and hangs the three list keywords off the array rather than off the element.
func applySchemaListKeys(sandbox *api.Sandbox, leaf *routeconf.Schema, props api.RouteBodyFieldProps) (*routeconf.Schema, error) {
	if !props.Array {
		if props.MinItems != "" || props.MaxItems != "" || props.UniqueItems {
			return nil, sandbox.Deps.Std.Errorf("--min-items, --max-items and --unique-items only apply to an --array property")
		}
		return leaf, nil
	}

	list := &routeconf.Schema{Type: "array", Items: leaf, UniqueItems: props.UniqueItems}
	if props.MinItems != "" {
		value, err := sandbox.Deps.Stringsdeps.Atoi(props.MinItems)
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("min-items must be an int, got %q", props.MinItems)
		}
		list.MinItems, list.HasMinItems = value, true
	}
	if props.MaxItems != "" {
		value, err := sandbox.Deps.Stringsdeps.Atoi(props.MaxItems)
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("max-items must be an int, got %q", props.MaxItems)
		}
		list.MaxItems, list.HasMaxItems = value, true
	}
	return list, nil
}

// RouteBodyFieldPropsOf reads one declared property back as the command line
// that would declare it — the inverse of RouteBodyPropertySchema, and what
// makes `set-body-field` an edit rather than a re-declaration: the keywords
// already there come back as props, the ones being changed are written over
// them, and the property is built again from the whole.
//
// Required is not read here: whether a property is demanded is its parent
// object's to say, not its own.
func RouteBodyFieldPropsOf(sandbox *api.Sandbox, schema *routeconf.Schema) api.RouteBodyFieldProps {
	props := api.RouteBodyFieldProps{Type: "string"}
	if schema == nil {
		return props
	}

	leaf := schema
	if schema.Type == "array" {
		props.Array = true
		props.UniqueItems = schema.UniqueItems
		if schema.HasMinItems {
			props.MinItems = sandbox.Deps.Stringsdeps.FormatInt(int64(schema.MinItems), 10)
		}
		if schema.HasMaxItems {
			props.MaxItems = sandbox.Deps.Stringsdeps.FormatInt(int64(schema.MaxItems), 10)
		}
		leaf = schema.Items
	}
	if leaf == nil {
		return props
	}

	props.Type = RouteDeclaredKind(leaf.Type)
	props.Nullable = leaf.Nullable
	props.Format = leaf.Format
	props.Pattern = leaf.Pattern
	props.Enum = leaf.Enum
	if leaf.HasConst {
		props.Const = leaf.Const
	}
	if leaf.HasAdditionalProperties {
		props.AdditionalProperties = leaf.AdditionalProperties
		props.NoAdditionalProperties = !leaf.AdditionalProperties
	}

	if leaf.HasMinLength {
		props.Min = sandbox.Deps.Stringsdeps.FormatInt(int64(leaf.MinLength), 10)
	}
	if leaf.HasMaxLength {
		props.Max = sandbox.Deps.Stringsdeps.FormatInt(int64(leaf.MaxLength), 10)
	}
	if leaf.HasMinimum {
		props.Min = RouteBoundText(sandbox, leaf.Minimum)
	}
	if leaf.HasMaximum {
		props.Max = RouteBoundText(sandbox, leaf.Maximum)
	}
	if leaf.HasExclusiveMinimum {
		props.ExclusiveMin = RouteBoundText(sandbox, leaf.ExclusiveMinimum)
	}
	if leaf.HasExclusiveMaximum {
		props.ExclusiveMax = RouteBoundText(sandbox, leaf.ExclusiveMaximum)
	}

	return props
}

// RouteDeclaredKind maps a json-schema type back onto the --type a command
// line spells it as, the inverse of RouteSchemaType.
func RouteDeclaredKind(kind string) string {
	switch kind {
	case "integer":
		return "int"
	case "number":
		return "float"
	case "boolean":
		return "boolean"
	case "object":
		return "object"
	default:
		return "string"
	}
}

// RouteBoundText writes a declared bound back as the literal a command line
// carries it as.
func RouteBoundText(sandbox *api.Sandbox, value float64) string {
	return sandbox.Deps.Stringsdeps.FormatFloat(value, 'g', -1, 64)
}

// SchemaDemands reports whether an object schema lists the named property in
// its `required` set. It is the one place a property's own declaration does
// not hold: whether a key has to be present is its parent's to say.
func SchemaDemands(parent *routeconf.Schema, name string) bool {
	if parent == nil {
		return false
	}
	for _, entry := range parent.Required {
		if entry == name {
			return true
		}
	}
	return false
}
