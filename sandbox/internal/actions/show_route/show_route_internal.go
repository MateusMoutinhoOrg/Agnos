package show_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// The indents the tree is drawn with. A route is three lists deep at most —
// the paths, the parameters and a body schema nested inside itself — so
// the depth is carried as a prefix rather than drawn with box characters,
// which keeps a line the same whether it is copied into a terminal or a doc.
const (
	branch = "  "
	level  = "  "
)

// ShowRouteInternal reads one route.yaml and renders it as the lines of a
// tree: the request line the route answers, then each place the declaration
// holds something — its paths, its parameters and the json-schema of its
// body, property by property.
//
// It is the one command of the route surface that writes nothing. A route is
// declared by several editors, each of which prints the file it changed and
// nothing else, so the declaration as a whole was only ever readable as yaml.
// This is that declaration said the way the editors talk about it.
func ShowRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, route string) ([]string, error) {
	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return nil, err
	}

	lines := []string{
		sandbox.Deps.Std.Sprintf("%s %s", sandbox.Deps.Stringsdeps.Join(conf.Methods, ","), conf.Pattern()),
	}
	lines = append(lines, routeHead(sandbox, conf)...)
	lines = append(lines, pathLines(sandbox, conf)...)
	lines = append(lines, parameterLines(sandbox, conf)...)
	lines = append(lines, bodyLines(sandbox, conf)...)

	return lines, nil
}

// routeHead is what the route says about itself: the package it is declared
// in, the heading it is listed under, its one-line help, the rung of the chain
// it runs on, and whether it is kept off the listings.
func routeHead(sandbox *api.Sandbox, conf *routeconf.RouteConf) []string {
	lines := []string{}

	if conf.Help != "" {
		lines = append(lines, branch+conf.Help)
	}
	if conf.Category != "" {
		lines = append(lines, sandbox.Deps.Std.Sprintf("%scategory  %s", branch, conf.Category))
	}
	lines = append(lines, sandbox.Deps.Std.Sprintf("%spriority  %d", branch, conf.Priority))
	lines = append(lines, sandbox.Deps.Std.Sprintf("%sresponse  %s", branch, conf.ResponseType))
	if conf.Hidden {
		lines = append(lines, branch+"hidden")
	}

	return lines
}

// pathLines is the route's paths, one line each: the Entries field it binds,
// the segments it reads and the trigger they have to match.
func pathLines(sandbox *api.Sandbox, conf *routeconf.RouteConf) []string {
	if len(conf.Paths) == 0 {
		return []string{}
	}

	lines := []string{"", "paths"}
	for _, path := range conf.Paths {
		text := sandbox.Deps.Std.Sprintf("%-20s segments %d..%d", path.Id, path.Start, path.End)

		notes := []string{}
		if path.Trigger.Exists {
			notes = append(notes, sandbox.Deps.Std.Sprintf("%s %q", path.Trigger.Type, path.Trigger.Value))
		}
		if path.Description != "" {
			notes = append(notes, path.Description)
		}
		lines = append(lines, branch+withNotes(sandbox, text, notes))
	}

	return lines
}

// parameterLines is the route's parameters, one line each: the Entries field
// it binds and the Go type it binds to, then where it is read from and every
// rule the dispatch holds a request to before the handler runs.
func parameterLines(sandbox *api.Sandbox, conf *routeconf.RouteConf) []string {
	if len(conf.Parameters) == 0 {
		return []string{}
	}

	lines := []string{"", "parameters"}
	for _, parameter := range conf.Parameters {
		text := sandbox.Deps.Std.Sprintf("%-20s %s", parameter.Id, parameter.Type)

		notes := []string{sandbox.Deps.Std.Sprintf("%q from %s", parameter.Key, sandbox.Deps.Stringsdeps.Join(parameter.Fonts, ", "))}
		if parameter.Required {
			notes = append(notes, "required")
		}
		if parameter.HasDefault {
			notes = append(notes, sandbox.Deps.Std.Sprintf("default %s", parameter.Default))
		}
		if parameter.Trigger.Exists {
			notes = append(notes, sandbox.Deps.Std.Sprintf("%s %q", parameter.Trigger.Type, parameter.Trigger.Value))
		}
		if parameter.Description != "" {
			notes = append(notes, parameter.Description)
		}
		lines = append(lines, branch+withNotes(sandbox, text, notes))
	}

	return lines
}

// withNotes joins one line's text and its notes, the notes said after two
// spaces and separated by commas.
func withNotes(sandbox *api.Sandbox, text string, notes []string) string {
	if len(notes) == 0 {
		return text
	}
	return sandbox.Deps.Std.Sprintf("%s  %s", text, sandbox.Deps.Stringsdeps.Join(notes, ", "))
}

// bodyLines is the request body: the envelope the dispatch settles before a
// handler runs, and the json-schema ReadBody holds the document to.
func bodyLines(sandbox *api.Sandbox, conf *routeconf.RouteConf) []string {
	if conf.Body.Type == routeconf.BodyNone {
		return []string{"", "body      none"}
	}

	head := sandbox.Deps.Std.Sprintf("body      %s", conf.Body.Type)
	if conf.Body.Required {
		head += ", required"
	}
	if conf.Body.ContentType != "" {
		head += ", " + conf.Body.ContentType
	}
	if conf.Body.MaxBytes > 0 {
		head += sandbox.Deps.Std.Sprintf(", at most %s bytes", sandbox.Deps.Stringsdeps.FormatInt(int64(conf.Body.MaxBytes), 10))
	}

	lines := []string{"", head}
	if conf.Body.Schema == nil {
		return lines
	}

	return append(lines, schemaLines(sandbox, conf.Body.Schema, branch)...)
}

// schemaLines walks one object schema, one line per property and one deeper
// level per object it holds. An array of objects is walked through its items,
// so a list of them reads like the object it is a list of.
func schemaLines(sandbox *api.Sandbox, schema *routeconf.Schema, indent string) []string {
	lines := []string{}

	for _, property := range schema.Properties {
		lines = append(lines, indent+propertyText(sandbox, schema, property))

		object := utils.SchemaObjectOf(property.Schema)
		if object != nil && object.Type == "object" {
			lines = append(lines, schemaLines(sandbox, object, indent+level)...)
		}
	}

	return lines
}

// propertyText is one property of the schema on one line: the json type it
// accepts, whether its object demands it, and every keyword declared on it.
func propertyText(sandbox *api.Sandbox, parent *routeconf.Schema, property routeconf.SchemaProperty) string {
	text := sandbox.Deps.Std.Sprintf("%-20s %s", property.Name, schemaTypeText(property.Schema))

	notes := []string{}
	if utils.SchemaDemands(parent, property.Name) {
		notes = append(notes, "required")
	}
	notes = append(notes, keywordNotes(sandbox, property.Schema)...)

	if len(notes) == 0 {
		return text
	}
	return sandbox.Deps.Std.Sprintf("%s  %s", text, sandbox.Deps.Stringsdeps.Join(notes, ", "))
}

// schemaTypeText is a property's json type, an array said as a list of the
// type it holds.
func schemaTypeText(schema *routeconf.Schema) string {
	object := utils.SchemaObjectOf(schema)
	if object == nil {
		return schema.Type
	}
	if schema.Type == "array" {
		return "[]" + object.Type
	}
	return schema.Type
}

// keywordNotes is every keyword one property declares, in the order
// add-body-field's flags list them. The array keywords are the list's own and
// the rest the element's, which is the one thing reading the yaml does not
// make obvious.
func keywordNotes(sandbox *api.Sandbox, schema *routeconf.Schema) []string {
	notes := []string{}

	if schema.Type == "array" {
		if schema.HasMinItems {
			notes = append(notes, sandbox.Deps.Std.Sprintf("min-items %d", schema.MinItems))
		}
		if schema.HasMaxItems {
			notes = append(notes, sandbox.Deps.Std.Sprintf("max-items %d", schema.MaxItems))
		}
		if schema.UniqueItems {
			notes = append(notes, "unique-items")
		}
	}

	leaf := utils.SchemaObjectOf(schema)
	if leaf == nil {
		return notes
	}

	if leaf.Nullable {
		notes = append(notes, "nullable")
	}
	if leaf.Format != "" {
		notes = append(notes, "format "+leaf.Format)
	}
	if leaf.Pattern != "" {
		notes = append(notes, "pattern "+leaf.Pattern)
	}
	if leaf.HasConst {
		notes = append(notes, "const "+leaf.Const)
	}
	if len(leaf.Enum) > 0 {
		notes = append(notes, "one of "+sandbox.Deps.Stringsdeps.Join(leaf.Enum, "|"))
	}
	if leaf.HasMinLength {
		notes = append(notes, sandbox.Deps.Std.Sprintf("min-length %d", leaf.MinLength))
	}
	if leaf.HasMaxLength {
		notes = append(notes, sandbox.Deps.Std.Sprintf("max-length %d", leaf.MaxLength))
	}
	if leaf.HasMinimum {
		notes = append(notes, "min "+utils.RouteBoundText(sandbox, leaf.Minimum))
	}
	if leaf.HasMaximum {
		notes = append(notes, "max "+utils.RouteBoundText(sandbox, leaf.Maximum))
	}
	if leaf.HasExclusiveMinimum {
		notes = append(notes, "over "+utils.RouteBoundText(sandbox, leaf.ExclusiveMinimum))
	}
	if leaf.HasExclusiveMaximum {
		notes = append(notes, "under "+utils.RouteBoundText(sandbox, leaf.ExclusiveMaximum))
	}
	if leaf.HasAdditionalProperties && leaf.AdditionalProperties {
		notes = append(notes, "undeclared keys accepted")
	}
	if leaf.HasAdditionalProperties && !leaf.AdditionalProperties {
		notes = append(notes, "undeclared keys refused")
	}

	return notes
}
