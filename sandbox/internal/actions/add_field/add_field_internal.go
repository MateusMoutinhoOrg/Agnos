package add_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddFieldInternal parses the target route's route.yaml, adds the new field to
// the origin --in names (refusing a duplicate) and writes the file back. The
// four origins differ only in where the entry lands, which is why one action
// covers them all.
func AddFieldInternal(deps *deps.Deps, io *smartio.SmartIO, props api.RouteFieldProps) error {
	in, err := utils.RouteFieldOrigin(deps, props.In)
	if err != nil {
		return err
	}

	conf, err := utils.LoadRouteConf(deps, io, props.Route)
	if err != nil {
		return err
	}

	deps.Std.Log("add-field adding %s to %s \n", in, utils.RouteConfPath(deps, props.Route))

	switch in {
	case utils.RouteFieldInPath:
		err = addPathSegment(deps, conf, props)
	case utils.RouteFieldInHeader:
		err = addHeader(deps, conf, props)
	case utils.RouteFieldInQuery:
		err = addParam(deps, conf, props)
	default:
		err = addBodyProperty(deps, conf, props)
	}
	if err != nil {
		return err
	}

	return utils.SaveRouteConf(deps, io, props.Route, conf)
}

// addPathSegment appends one segment to `paths`: a trigger when --identifier
// is given, a capture when --name is. A trigger is normalized to start with
// "/", so "users" and "/users" are the same request.
func addPathSegment(deps *deps.Deps, conf *routeconf.RouteConf, props api.RouteFieldProps) error {
	position, err := utils.CheckRoutePosition(deps, "segment", props.Position, len(conf.Paths))
	if err != nil {
		return err
	}

	if deps.Stringsdeps.TrimSpace(props.Identifier) != "" {
		identifier, err := utils.RouteIdentifierSegment(deps, props.Identifier)
		if err != nil {
			return err
		}
		conf.Paths = utils.InsertRouteSegment(conf.Paths, routeconf.Segment{Identifier: identifier}, position)
		return nil
	}

	field, err := utils.NewRouteField(deps, props, utils.RouteFieldInPath)
	if err != nil {
		return err
	}
	if utils.FindRouteSegment(deps, conf.Paths, field.Key) >= 0 {
		return deps.Std.Errorf("route %q already captures a segment named %q", props.Route, field.Key)
	}

	conf.Paths = utils.InsertRouteSegment(conf.Paths, routeconf.Segment{Field: &field}, position)
	return nil
}

// addHeader appends one declared request header.
func addHeader(deps *deps.Deps, conf *routeconf.RouteConf, props api.RouteFieldProps) error {
	field, err := utils.NewRouteField(deps, props, utils.RouteFieldInHeader)
	if err != nil {
		return err
	}
	if utils.FindRouteField(deps, conf.Headers, field.Key) >= 0 {
		return deps.Std.Errorf("route %q already declares a header named %q", props.Route, field.Key)
	}

	position, err := utils.CheckRoutePosition(deps, "header", props.Position, len(conf.Headers))
	if err != nil {
		return err
	}

	conf.Headers = utils.InsertRouteField(conf.Headers, field, position)
	return nil
}

// addParam appends one declared query-string parameter.
func addParam(deps *deps.Deps, conf *routeconf.RouteConf, props api.RouteFieldProps) error {
	field, err := utils.NewRouteField(deps, props, utils.RouteFieldInQuery)
	if err != nil {
		return err
	}
	if utils.FindRouteField(deps, conf.Params, field.Key) >= 0 {
		return deps.Std.Errorf("route %q already declares a query parameter named %q", props.Route, field.Key)
	}

	position, err := utils.CheckRoutePosition(deps, "query parameter", props.Position, len(conf.Params))
	if err != nil {
		return err
	}

	conf.Params = utils.InsertRouteField(conf.Params, field, position)
	return nil
}

// addBodyProperty declares one property of the body's json-schema. The name is
// a dotted path ("address.city"), and the objects it passes through are
// created as they are needed. A route that declared no body becomes a json one
// here: declaring a property is what says it takes a body at all.
func addBodyProperty(deps *deps.Deps, conf *routeconf.RouteConf, props api.RouteFieldProps) error {
	name := utils.RouteFieldName(deps, props.Name)
	if name == "" {
		return deps.Std.Errorf("a body property needs a name")
	}

	if conf.Body.Type == routeconf.BodyNone {
		conf.Body.Type = "json"
		conf.Body.ContentType = routeconf.DefaultJsonContentType
	}
	if conf.Body.Type != "json" {
		return deps.Std.Errorf("route %q declares a %q body, which carries no json-schema", props.Route, conf.Body.Type)
	}
	if conf.Body.Schema == nil {
		conf.Body.Schema = &routeconf.Schema{Type: "object"}
		conf.Body.HasSchema = true
	}

	parts := deps.Stringsdeps.Split(name, ".")
	parent := conf.Body.Schema

	for _, segment := range parts[:len(parts)-1] {
		child := findProperty(parent, segment)
		if child == nil {
			child = &routeconf.Schema{Type: "object"}
			insertProperty(parent, segment, child)
		}
		if child.Type != "object" {
			return deps.Std.Errorf("body property %q is not an object, so nothing can be declared under it", segment)
		}
		parent = child
	}

	leaf := parts[len(parts)-1]
	if findProperty(parent, leaf) != nil {
		return deps.Std.Errorf("route %q already declares a body property named %q", props.Route, name)
	}

	schema, err := propertySchema(deps, props)
	if err != nil {
		return err
	}
	insertProperty(parent, leaf, schema)

	if props.Required {
		parent.Required = utils.AppendUnique(parent.Required, []string{leaf})
		deps.Sortdeps.Strings(parent.Required)
	}

	return nil
}

// propertySchema turns the raw values typed on the command line into one
// json-schema node: the type, the bounds it accepts for that type, and the
// two string refinements --format and --pattern.
func propertySchema(deps *deps.Deps, props api.RouteFieldProps) (*routeconf.Schema, error) {
	kind, ok := utils.FieldType(deps, props.Type)
	if !ok {
		return nil, deps.Std.Errorf("unknown type %q (use string, boolean, int or float)", props.Type)
	}

	leaf := &routeconf.Schema{
		Type:    schemaType(kind),
		Format:  deps.Stringsdeps.TrimSpace(props.Format),
		Pattern: deps.Stringsdeps.TrimSpace(props.Pattern),
	}

	if (leaf.Format != "" || leaf.Pattern != "") && kind != "string" {
		return nil, deps.Std.Errorf("--format and --pattern only apply to a string property")
	}

	if err := applyBounds(deps, leaf, kind, props); err != nil {
		return nil, err
	}

	if !props.Array {
		return leaf, nil
	}
	return &routeconf.Schema{Type: "array", Items: leaf}, nil
}

// applyBounds spells --min/--max as the keywords the property's type uses:
// minimum/maximum for a number, minLength/maxLength for text.
func applyBounds(deps *deps.Deps, schema *routeconf.Schema, kind string, props api.RouteFieldProps) error {
	if props.Min == "" && props.Max == "" {
		return nil
	}
	if kind == "boolean" {
		return deps.Std.Errorf("min/max do not apply to a boolean property")
	}

	if kind == "string" {
		if props.Min != "" {
			value, err := deps.Stringsdeps.Atoi(props.Min)
			if err != nil {
				return deps.Std.Errorf("min must be an int, got %q", props.Min)
			}
			schema.MinLength, schema.HasMinLength = value, true
		}
		if props.Max != "" {
			value, err := deps.Stringsdeps.Atoi(props.Max)
			if err != nil {
				return deps.Std.Errorf("max must be an int, got %q", props.Max)
			}
			schema.MaxLength, schema.HasMaxLength = value, true
		}
		return nil
	}

	if props.Min != "" {
		value, err := utils.RouteParseBound(deps, kind, "min", props.Min)
		if err != nil {
			return err
		}
		schema.Minimum, schema.HasMinimum = value, true
	}
	if props.Max != "" {
		value, err := utils.RouteParseBound(deps, kind, "max", props.Max)
		if err != nil {
			return err
		}
		schema.Maximum, schema.HasMaximum = value, true
	}
	return nil
}

// schemaType maps a declared field type onto its json-schema spelling.
func schemaType(kind string) string {
	switch kind {
	case "int":
		return "integer"
	case "float":
		return "number"
	case "boolean":
		return "boolean"
	default:
		return "string"
	}
}

// findProperty returns the named property of an object schema, or nil.
func findProperty(schema *routeconf.Schema, name string) *routeconf.Schema {
	for _, property := range schema.Properties {
		if property.Name == name {
			return property.Schema
		}
	}
	return nil
}

// insertProperty adds one property in alphabetical position, which is the
// order the declaration is read back in — so the generated struct and the
// rendered route.yaml are the same on every build.
func insertProperty(parent *routeconf.Schema, name string, schema *routeconf.Schema) {
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
