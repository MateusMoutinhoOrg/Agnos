package show_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// The indents the tree is drawn with. A route is four lists deep at most — the
// path, the headers, the params and a body schema nested inside itself — so
// the depth is carried as a prefix rather than drawn with box characters,
// which keeps a line the same whether it is copied into a terminal or a doc.
const (
	branch = "  "
	level  = "  "
)

// ShowRouteInternal reads one route.yaml and renders it as the lines of a
// tree: the request line the route answers, then each place the declaration
// holds something — the segments of its path, its headers, its query
// parameters and the json-schema of its body, property by property.
//
// It is the one command of the route surface that writes nothing. A route is
// declared by six editors, each of which prints the file it changed and
// nothing else, so the declaration as a whole was only ever readable as yaml.
// This is that declaration said the way the editors talk about it.
func ShowRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, route string) ([]string, error) {
	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return nil, err
	}

	lines := []string{
		sandbox.Deps.Std.Sprintf("%s %s", conf.Method, conf.Pattern()),
	}
	lines = append(lines, routeHead(sandbox, conf)...)
	lines = append(lines, pathLines(sandbox, conf)...)
	lines = append(lines, fieldLines(sandbox, "headers", conf.Headers)...)
	lines = append(lines, fieldLines(sandbox, "params", conf.Params)...)
	lines = append(lines, bodyLines(sandbox, conf)...)

	return lines, nil
}

// routeHead is what the route says about itself: the package it is declared
// in, the heading it is listed under, its one-line help, and whether it is
// kept off the listings.
func routeHead(sandbox *api.Sandbox, conf *routeconf.RouteConf) []string {
	lines := []string{}

	if conf.Help != "" {
		lines = append(lines, branch+conf.Help)
	}
	if conf.Category != "" {
		lines = append(lines, sandbox.Deps.Std.Sprintf("%scategory  %s", branch, conf.Category))
	}
	if conf.Hidden {
		lines = append(lines, branch+"hidden")
	}

	return lines
}

// pathLines is the route's path, segment by segment in the order the URL
// spells them: a trigger by the literal it matches, a capture by the field it
// binds.
func pathLines(sandbox *api.Sandbox, conf *routeconf.RouteConf) []string {
	if len(conf.Paths) == 0 {
		return []string{}
	}

	lines := []string{"", "path"}
	for _, segment := range conf.Paths {
		if segment.Field == nil {
			lines = append(lines, sandbox.Deps.Std.Sprintf("%s%s", branch, segment.Identifier))
			continue
		}
		lines = append(lines, branch+fieldText(sandbox, *segment.Field, "{"+segment.Field.Key+"}"))
	}

	return lines
}

// fieldLines is one of the two lists that read a value off the request line,
// under the heading its editors name it by.
func fieldLines(sandbox *api.Sandbox, heading string, fields []routeconf.Field) []string {
	if len(fields) == 0 {
		return []string{}
	}

	lines := []string{"", heading}
	for _, field := range fields {
		lines = append(lines, branch+fieldText(sandbox, field, field.Key))
	}

	return lines
}

// fieldText is one header, query parameter or captured segment on one line:
// what it is called, the Go type it binds to, and every rule the dispatch
// holds a request to before the handler runs.
func fieldText(sandbox *api.Sandbox, field routeconf.Field, name string) string {
	text := sandbox.Deps.Std.Sprintf("%-20s %s", name, typeText(field))

	notes := []string{}
	if field.Required {
		notes = append(notes, "required")
	}
	if field.HasDefault {
		notes = append(notes, sandbox.Deps.Std.Sprintf("default %s", field.Default))
	}
	if field.HasMin {
		notes = append(notes, sandbox.Deps.Std.Sprintf("min %s", utils.RouteBoundText(sandbox, field.Min)))
	}
	if field.HasMax {
		notes = append(notes, sandbox.Deps.Std.Sprintf("max %s", utils.RouteBoundText(sandbox, field.Max)))
	}
	if field.Description != "" {
		notes = append(notes, field.Description)
	}

	if len(notes) == 0 {
		return text
	}
	return sandbox.Deps.Std.Sprintf("%s  %s", text, sandbox.Deps.Stringsdeps.Join(notes, ", "))
}

// typeText is a field's declared type, said as the Go value it binds to.
func typeText(field routeconf.Field) string {
	if field.Array {
		return "[]" + field.Type
	}
	return field.Type
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
