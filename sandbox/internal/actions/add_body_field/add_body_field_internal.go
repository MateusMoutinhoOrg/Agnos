package add_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddBodyFieldInternal parses the target route's route.yaml, declares one
// property of the body's json-schema at the dotted path props.Name and writes
// the file back. A route that declared no body becomes a json one here:
// declaring a property is what says it takes a body at all.
func AddBodyFieldInternal(deps *deps.Deps, io *smartio.SmartIO, props api.RouteBodyFieldProps) error {
	conf, err := utils.LoadRouteConf(deps, io, props.Route)
	if err != nil {
		return err
	}

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

	parts := utils.SplitSchemaPath(deps, name)
	parent, err := walkTo(deps, conf.Body.Schema, parts[:len(parts)-1])
	if err != nil {
		return err
	}

	leaf := parts[len(parts)-1]
	if utils.SchemaPropertyOf(parent, leaf) != nil {
		return deps.Std.Errorf("route %q already declares a body property named %q", props.Route, name)
	}

	schema, err := propertySchema(deps, props)
	if err != nil {
		return err
	}

	deps.Std.Log("add-body-field adding %s to %s \n", name, utils.RouteConfPath(deps, props.Route))

	utils.InsertSchemaProperty(parent, leaf, schema)
	if props.Required {
		parent.Required = utils.AppendUnique(parent.Required, []string{leaf})
		deps.Sortdeps.Strings(parent.Required)
	}

	return utils.SaveRouteConf(deps, io, props.Route, conf)
}

// walkTo descends the dotted path to the object the leaf is declared in,
// creating the objects it passes through. An array of objects is entered
// through its `items`, so "tags.name" reads the same whether tags is an object
// or a list of them.
func walkTo(deps *deps.Deps, root *routeconf.Schema, segments []string) (*routeconf.Schema, error) {
	parent := root
	for _, segment := range segments {
		property := utils.SchemaPropertyOf(parent, segment)
		if property == nil {
			child := &routeconf.Schema{Type: "object"}
			utils.InsertSchemaProperty(parent, segment, child)
			parent = child
			continue
		}

		child := utils.SchemaObjectOf(property)
		if child == nil || child.Type != "object" {
			return nil, deps.Std.Errorf("body property %q is not an object, so nothing can be declared under it", segment)
		}
		parent = child
	}
	return parent, nil
}

// propertySchema turns the raw values typed on the command line into one
// json-schema node. --array wraps the whole of it, so every keyword but the
// three list ones describes the element and not the list.
func propertySchema(deps *deps.Deps, props api.RouteBodyFieldProps) (*routeconf.Schema, error) {
	kind, err := utils.RouteSchemaKind(deps, props.Type)
	if err != nil {
		return nil, err
	}
	format, err := utils.RouteSchemaFormat(deps, props.Format)
	if err != nil {
		return nil, err
	}

	leaf := &routeconf.Schema{
		Type:     utils.RouteSchemaType(kind),
		Nullable: props.Nullable,
		Format:   format,
		Pattern:  deps.Stringsdeps.TrimSpace(props.Pattern),
		Enum:     props.Enum,
	}
	if (leaf.Format != "" || leaf.Pattern != "") && kind != "string" {
		return nil, deps.Std.Errorf("--format and --pattern only apply to a string property")
	}
	if value := deps.Stringsdeps.TrimSpace(props.Const); value != "" {
		leaf.Const, leaf.HasConst = value, true
	}

	if err := applyBounds(deps, leaf, kind, props); err != nil {
		return nil, err
	}
	if err := applyObjectKeys(deps, leaf, kind, props); err != nil {
		return nil, err
	}

	return applyListKeys(deps, leaf, props)
}

// applyBounds spells --min/--max as the keywords the property's type uses:
// minimum/maximum for a number, minLength/maxLength for text. The two
// exclusive bounds are a number's alone.
func applyBounds(deps *deps.Deps, schema *routeconf.Schema, kind string, props api.RouteBodyFieldProps) error {
	numeric := kind == "int" || kind == "float"

	if props.ExclusiveMin != "" || props.ExclusiveMax != "" {
		if !numeric {
			return deps.Std.Errorf("--exclusive-min and --exclusive-max only apply to an int or float property")
		}
		if props.ExclusiveMin != "" {
			value, err := utils.RouteParseBound(deps, kind, "exclusive-min", props.ExclusiveMin)
			if err != nil {
				return err
			}
			schema.ExclusiveMinimum, schema.HasExclusiveMinimum = value, true
		}
		if props.ExclusiveMax != "" {
			value, err := utils.RouteParseBound(deps, kind, "exclusive-max", props.ExclusiveMax)
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
		return deps.Std.Errorf("min/max do not apply to a %s property", kind)
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

// applyObjectKeys settles whether an object property accepts keys it does not
// declare. It is the one switch that has to be written even when it is false,
// so the schema tells "undeclared keys are refused" from "nothing was said".
func applyObjectKeys(deps *deps.Deps, schema *routeconf.Schema, kind string, props api.RouteBodyFieldProps) error {
	if !props.AdditionalProperties && !props.NoAdditionalProperties {
		return nil
	}
	if props.AdditionalProperties && props.NoAdditionalProperties {
		return deps.Std.Errorf("--additional-properties and --no-additional-properties are mutually exclusive")
	}
	if kind != "object" {
		return deps.Std.Errorf("--additional-properties only applies to an object property")
	}

	schema.AdditionalProperties = props.AdditionalProperties
	schema.HasAdditionalProperties = true
	return nil
}

// applyListKeys wraps the element schema in the array --array asks for, and
// hangs the three list keywords off the array rather than off the element.
func applyListKeys(deps *deps.Deps, leaf *routeconf.Schema, props api.RouteBodyFieldProps) (*routeconf.Schema, error) {
	if !props.Array {
		if props.MinItems != "" || props.MaxItems != "" || props.UniqueItems {
			return nil, deps.Std.Errorf("--min-items, --max-items and --unique-items only apply to an --array property")
		}
		return leaf, nil
	}

	list := &routeconf.Schema{Type: "array", Items: leaf, UniqueItems: props.UniqueItems}
	if props.MinItems != "" {
		value, err := deps.Stringsdeps.Atoi(props.MinItems)
		if err != nil {
			return nil, deps.Std.Errorf("min-items must be an int, got %q", props.MinItems)
		}
		list.MinItems, list.HasMinItems = value, true
	}
	if props.MaxItems != "" {
		value, err := deps.Stringsdeps.Atoi(props.MaxItems)
		if err != nil {
			return nil, deps.Std.Errorf("max-items must be an int, got %q", props.MaxItems)
		}
		list.MaxItems, list.HasMaxItems = value, true
	}
	return list, nil
}
