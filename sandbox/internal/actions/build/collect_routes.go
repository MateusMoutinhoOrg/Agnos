package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// routesDir holds one declared route per sub-directory, the server layer's
// mirror of sandbox/internal/commands.
const routesDir = "sandbox/internal/routes"

// CollectRoutes reads every sandbox/internal/routes/<name>/route.yaml and
// returns one rich data map per route, for the {{range .Routes}} loops in the
// generated sandbox/internal/server/servermain.go and for the per-route
// entries.go. It is the server layer's CollectCommands.
//
// The list comes back ordered for matching, not in listing order: the most
// specific route first, so a route fixing "/" can never swallow "/home".
func CollectRoutes(deps *deps.Deps, io *smartio.SmartIO) ([]map[string]any, error) {
	var routes []map[string]any

	for _, dir := range io.ListDirs(routesDir) {
		name := lastSegmentOf(deps, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile(routesDir + "/" + name + "/route.yaml")
		if err != nil {
			continue
		}

		conf, err := routeconf.New(deps, string(content))
		if err != nil {
			return nil, deps.Std.Errorf("routes/%s/route.yaml: %w", name, err)
		}

		if err := checkRouteSchema(deps, name, conf); err != nil {
			return nil, err
		}

		if err := checkRoutePaths(deps, name, conf); err != nil {
			return nil, err
		}

		routes = append(routes, routeData(deps, name, conf))
	}

	sortRoutes(deps, routes)
	return routes, nil
}

// sortRoutes puts the routes in the order the dispatch tests them: the route
// fixing the most literal segments first, then the one whose literals spell
// the most characters, then the one of fixed length before the one taking the
// rest of the path, then by pattern for a stable tie-break. The ordering is
// the collector's, so servermain.go's template only has to range in order.
func sortRoutes(deps *deps.Deps, routes []map[string]any) {
	deps.Sortdeps.SliceStable(routes, func(i int, j int) bool {
		left, right := routes[i], routes[j]
		if left["IdentifierCount"] != right["IdentifierCount"] {
			return left["IdentifierCount"].(int) > right["IdentifierCount"].(int)
		}
		if left["IdentifierLen"] != right["IdentifierLen"] {
			return left["IdentifierLen"].(int) > right["IdentifierLen"].(int)
		}
		if left["HasRest"] != right["HasRest"] {
			return !left["HasRest"].(bool)
		}
		return left["Pattern"].(string) < right["Pattern"].(string)
	})
}

// checkRouteSchema refuses a json-schema that reaches outside the declared
// subset. Silently dropping $ref or oneOf would generate a validator that
// accepts what the declaration means to reject.
func checkRouteSchema(deps *deps.Deps, name string, conf *routeconf.RouteConf) error {
	for _, key := range schemaUnknownKeys(conf.Body.Schema) {
		return deps.Std.Errorf("routes/%s/route.yaml: json-schema key %q is outside the supported subset", name, key)
	}
	return nil
}

// schemaUnknownKeys walks a schema tree and gathers every key outside the
// subset, whatever depth it sits at.
func schemaUnknownKeys(schema *routeconf.Schema) []string {
	if schema == nil {
		return nil
	}
	unknown := schema.Unknown
	unknown = append(unknown, schemaUnknownKeys(schema.Items)...)
	for _, property := range schema.Properties {
		unknown = append(unknown, schemaUnknownKeys(property.Schema)...)
	}
	return unknown
}

// checkRoutePaths refuses a `paths` a matcher cannot be generated from: the
// segment taking the rest of the path only reads as a suffix if it is the last
// one, and its Entries field is a []T, which no header or query parameter of
// the same name could fill.
func checkRoutePaths(deps *deps.Deps, name string, conf *routeconf.RouteConf) error {
	at := utils.RouteRestIndex(conf.Paths)
	if at < 0 {
		return nil
	}

	if at != len(conf.Paths)-1 {
		return deps.Std.Errorf("routes/%s/route.yaml: the captured segment %q is an array before the end of `paths`; only the last segment takes the rest of the path", name, conf.Paths[at].Field.Key)
	}

	key := conf.Paths[at].Field.Key
	for _, origin := range []struct {
		label  string
		fields []routeconf.Field
	}{{"header", conf.Headers}, {"query parameter", conf.Params}} {
		if utils.FindRouteField(deps, origin.fields, key) >= 0 {
			return deps.Std.Errorf("routes/%s/route.yaml: %q is both the segment taking the rest of the path and a declared %s; the two cannot fill one Entries field", name, key, origin.label)
		}
	}

	return nil
}

// arityOp is the comparison match<Route> refuses a path with. A route of fixed
// length wants exactly the segments it declares; one taking the rest of the
// path wants strictly more, because that last capture is required like any
// other.
func arityOp(rest_index int) string {
	if rest_index < 0 {
		return "!="
	}
	return "<="
}

// routeData is one route as both templates read it.
func routeData(deps *deps.Deps, name string, conf *routeconf.RouteConf) map[string]any {
	bindings, parts, count, rest_index := routeBindings(deps, conf)

	return map[string]any{
		"Name":            name,
		"GoName":          exportedName(deps, name),
		"Method":          conf.Method,
		"Trigger":         routeTrigger(conf),
		"Pattern":         conf.Pattern(),
		"IdentifierCount": conf.IdentifierCount(),
		"IdentifierLen":   conf.IdentifierLen(),
		"SegmentCount":    count,
		"HasRest":         rest_index >= 0,
		"ArityOp":         arityOp(rest_index),
		"MatchParts":      parts,
		"Bindings":        bindings,
		"Category":        conf.Category,
		"Help":            conf.Help,
		"LongDescription": conf.LongDescription,
		"Examples":        conf.Examples,
		"Hidden":          conf.Hidden,
		"Body":            bodyData(deps, conf),
		"HasBody":         conf.Body.Type != routeconf.BodyNone,
		"SchemaJson":      conf.SchemaJson(),
		"BodyStructs":     bodyStructs(deps, conf),
	}
}

// routeTrigger is the first literal segment of a route — the one that names it
// in docs/Routes. A route made of captures alone falls back to its pattern.
func routeTrigger(conf *routeconf.RouteConf) string {
	for _, segment := range conf.Paths {
		if segment.Field == nil {
			return segment.Identifier
		}
	}
	return conf.Pattern()
}

// routeBindings resolves the route's whole input surface into one entry per
// generated Entries field, the URL segments the dispatch matches on, how many
// segments a matching path fixes, and the position the last capture takes the
// rest of the path from (-1 when the route fixes its whole length).
//
// The same name may be declared in more than one origin: the field is written
// once and filled by the first origin, in declaration order, that brings a
// value.
func routeBindings(deps *deps.Deps, conf *routeconf.RouteConf) ([]map[string]any, []map[string]any, int, int) {
	var bindings []map[string]any
	var parts []map[string]any
	index := map[string]int{}
	position := 0
	rest_index := -1

	for _, segment := range conf.Paths {
		if segment.Field == nil {
			// The root identifier "/" fixes the empty path: it names no
			// segment of its own, so it adds nothing to match on.
			literal := deps.Stringsdeps.Trim(segment.Identifier, "/")
			if literal == "" {
				continue
			}
			parts = append(parts, map[string]any{
				"Index":     position,
				"IsCapture": false,
				"Literal":   literal,
			})
			position++
			continue
		}
		if segment.Field.Array {
			// The segment taking the rest of the path fixes no position
			// of its own: it opens the suffix, so it adds nothing to
			// match on and closes the sequence.
			rest_index = position
			bindings = addBinding(deps, bindings, index, *segment.Field, "path", position)
			break
		}

		parts = append(parts, map[string]any{
			"Index":     position,
			"IsCapture": true,
			"Literal":   "",
		})
		bindings = addBinding(deps, bindings, index, *segment.Field, "path", position)
		position++
	}

	for _, field := range conf.Headers {
		bindings = addBinding(deps, bindings, index, field, "header", 0)
	}
	for _, field := range conf.Params {
		bindings = addBinding(deps, bindings, index, field, "query", 0)
	}

	return bindings, parts, position, rest_index
}

// addBinding records one declared origin against the Entries field its name
// resolves to, creating that field the first time the name is seen.
func addBinding(deps *deps.Deps, bindings []map[string]any, index map[string]int, field routeconf.Field, kind string, position int) []map[string]any {
	source := map[string]any{
		"Kind":   kind,
		"Index":  position,
		"Key":    field.Key,
		"Getter": sourceGetter(deps, kind, field.Key),
	}

	if at, seen := index[field.Key]; seen {
		binding := bindings[at]
		binding["Sources"] = append(binding["Sources"].([]map[string]any), source)
		if field.Required {
			binding["Required"] = true
		}
		if field.HasDefault && !binding["HasDefault"].(bool) {
			binding["HasDefault"] = true
			binding["DefaultLiteral"] = routeDefaultLiteral(deps, field.Type, field.Default)
		}
		setBindingArity(binding, len(binding["Sources"].([]map[string]any)))
		return bindings
	}

	index[field.Key] = len(bindings)
	binding := map[string]any{
		"Key":            field.Key,
		"GoField":        exportedName(deps, field.Key),
		"GoType":         goType(field.Type, field.Array),
		"Type":           field.Type,
		"IsArray":        field.Array,
		"Subject":        subjectLabel(kind),
		"ParseFunc":      routeParseFunc(field.Type),
		"Required":       field.Required,
		"HasDefault":     field.HasDefault,
		"DefaultLiteral": routeDefaultLiteral(deps, field.Type, field.Default),
		"Description":    field.Description,
		"Examples":       field.Examples,
		"Default":        field.Default,
		"MinLabel":       routeNumberLabel(deps, field.Type, field.Min, field.HasMin),
		"MaxLabel":       routeNumberLabel(deps, field.Type, field.Max, field.HasMax),
		"RangeCheck":     routeRangeCheck(deps, field, kind),
		"Sources":        []map[string]any{source},
	}
	binding["HasPath"] = kind == "path"
	binding["IsRest"] = kind == "path" && field.Array
	binding["RestIndex"] = position
	setBindingArity(binding, 1)
	return append(bindings, binding)
}

// setBindingArity settles the two flags the generated dispatch is shaped by.
// A binding fed by a captured path segment is always filled — the route only
// matched because that segment was there — so it needs no fallback, and one
// fed by a single origin needs no "was it filled?" flag at all.
func setBindingArity(binding map[string]any, sources int) {
	binding["MultiSource"] = sources > 1
	binding["NeedsFallback"] = (binding["Required"].(bool) || binding["HasDefault"].(bool)) && !binding["HasPath"].(bool)
	binding["NeedsFound"] = sources > 1 || binding["NeedsFallback"].(bool)
}

// sourceGetter is the expression the generated dispatch reads one origin's raw
// value with. A path capture is read by index instead, so it has none.
func sourceGetter(deps *deps.Deps, kind string, key string) string {
	switch kind {
	case "header":
		return "request.GetHeader(" + deps.Stringsdeps.Quote(key) + ")"
	case "query":
		return "request.GetQueryParam(" + deps.Stringsdeps.Quote(key) + ")"
	default:
		return ""
	}
}

// subjectLabel is how an error message names the origin a value came from.
func subjectLabel(kind string) string {
	switch kind {
	case "header":
		return "header"
	case "query":
		return "query parameter"
	default:
		return "path segment"
	}
}

// routeParseFunc is the generated helper that converts one raw request value
// into the field's Go type, answering 400 on its own when it cannot.
func routeParseFunc(kind string) string {
	switch kind {
	case "int":
		return "parseIntValue"
	case "float":
		return "parseFloatValue"
	case "boolean":
		return "parseBoolValue"
	default:
		return "parseStringValue"
	}
}

// routeNumberLabel renders a min/max bound as the literal it has in
// route.yaml ("" when the bound is unset).
func routeNumberLabel(deps *deps.Deps, kind string, value float64, has bool) string {
	if !has {
		return ""
	}
	if kind == "int" {
		return deps.Stringsdeps.FormatInt(int64(value), 10)
	}
	return deps.Stringsdeps.FormatFloat(value, 'g', -1, 64)
}

// routeDefaultLiteral is the Go literal a declared default is assigned as.
func routeDefaultLiteral(deps *deps.Deps, kind string, value string) string {
	switch kind {
	case "boolean":
		if value == "true" {
			return "true"
		}
		return "false"
	case "int", "float":
		if value == "" {
			return "0"
		}
		return value
	default:
		return deps.Stringsdeps.Quote(value)
	}
}

// routeRangeCheck emits the Go statements the generated dispatch runs, after a
// numeric field has been bound, to enforce its min/max bounds. It returns ""
// for fields that carry no bound (or are not int/float scalars). The body is
// indented one tab — the depth of the block it is spliced into.
func routeRangeCheck(deps *deps.Deps, field routeconf.Field, kind string) string {
	if field.Array || (field.Type != "int" && field.Type != "float") {
		return ""
	}
	if !field.HasMin && !field.HasMax {
		return ""
	}

	subject := subjectLabel(kind)
	goField := exportedName(deps, field.Key)

	b := ""
	// failOp is the comparison that means "out of range"; wantOp is what the
	// message tells the caller to satisfy.
	guard := func(failOp, wantOp, bound string) {
		b += deps.Std.Sprintf(
			"\tif entries.%s %s %s {\n"+
				"\t\trouteio.WriteError(deps, response, api.StatusBadRequest, %q, \"%s '%s' must be %s %s\")\n"+
				"\t\treturn\n"+
				"\t}\n",
			goField, failOp, bound, field.Key, subject, field.Key, wantOp, bound)
	}
	if field.HasMin {
		guard("<", ">=", routeNumberLabel(deps, field.Type, field.Min, true))
	}
	if field.HasMax {
		guard(">", "<=", routeNumberLabel(deps, field.Type, field.Max, true))
	}
	return deps.Stringsdeps.TrimRight(b, "\n")
}

// bodyData is the route's body declaration as the two templates read it.
func bodyData(deps *deps.Deps, conf *routeconf.RouteConf) map[string]any {
	body := conf.Body
	return map[string]any{
		"Type":         body.Type,
		"Required":     body.Required,
		"MaxBytes":     body.MaxBytes,
		"MaxBytesText": deps.Stringsdeps.FormatInt(int64(body.MaxBytes), 10),
		"ContentType":  body.ContentType,
		"HasSchema":    body.HasSchema,
		"IsRaw":        body.Type == "raw",
		"IsText":       body.Type == "text",
		"IsJson":       body.Type == "json",
		"IsObject":     body.HasSchema && body.Schema != nil && body.Schema.Type == "object",
		"GoType":       bodyGoType(body),
	}
}

// bodyGoType is what the generated ReadBody hands back: the bytes themselves,
// the text, the struct the schema describes, or — for a json body with no
// schema — the parsed document.
func bodyGoType(body routeconf.Body) string {
	switch body.Type {
	case "raw":
		return "[]byte"
	case "text":
		return "string"
	case "json":
		if body.HasSchema && body.Schema != nil && body.Schema.Type == "object" {
			return "Body"
		}
		return "*serializables.SerializibleObject"
	}
	return ""
}

// bodyStructs flattens the declared json-schema into the Go structs the
// generated entries.go declares: Body for the root object, Body<Path> for a
// nested object and Body<Path>Item for the object an array holds.
func bodyStructs(deps *deps.Deps, conf *routeconf.RouteConf) []map[string]any {
	body := conf.Body
	if body.Type != "json" || !body.HasSchema || body.Schema == nil || body.Schema.Type != "object" {
		return nil
	}
	return appendBodyStruct(deps, nil, "Body", body.Schema)
}

// appendBodyStruct emits one struct for schema and, depth first, one for every
// object nested under it. The parent is appended before its children so the
// generated file reads outside in.
func appendBodyStruct(deps *deps.Deps, structs []map[string]any, name string, schema *routeconf.Schema) []map[string]any {
	fields := make([]map[string]any, 0, len(schema.Properties))
	var nested []*routeconf.Schema
	var nested_names []string

	for _, property := range schema.Properties {
		go_field := exportedName(deps, property.Name)
		child := property.Schema
		field := map[string]any{
			"Key":     property.Name,
			"GoField": go_field,
			"Kind":    child.Type,
		}

		field["Reader"] = schemaReader(child.Type)

		switch child.Type {
		case "object":
			struct_name := name + go_field
			field["GoType"] = struct_name
			field["StructName"] = struct_name
			nested = append(nested, child)
			nested_names = append(nested_names, struct_name)
		case "array":
			item := child.Items
			if item != nil && item.Type == "object" {
				struct_name := name + go_field + "Item"
				field["GoType"] = "[]" + struct_name
				field["StructName"] = struct_name
				field["ItemKind"] = "object"
				field["ItemReader"] = ""
				nested = append(nested, item)
				nested_names = append(nested_names, struct_name)
				break
			}
			item_kind := ""
			if item != nil {
				item_kind = item.Type
			}
			field["GoType"] = "[]" + schemaGoType(item_kind)
			field["ItemKind"] = item_kind
			field["ItemReader"] = schemaItemReader(item_kind)
		default:
			field["GoType"] = schemaGoType(child.Type)
		}

		fields = append(fields, field)
	}

	structs = append(structs, map[string]any{
		"Name":   name,
		"Fields": fields,
	})

	for i, child := range nested {
		structs = appendBodyStruct(deps, structs, nested_names[i], child)
	}

	return structs
}

// schemaReader is the routeio reader the generated bind function pulls one
// scalar property with.
func schemaReader(kind string) string {
	switch kind {
	case "integer":
		return "ReadInt"
	case "number":
		return "ReadFloat"
	case "boolean":
		return "ReadBool"
	default:
		return "ReadString"
	}
}

// schemaItemReader is schemaReader for one item of an array of scalars.
func schemaItemReader(kind string) string {
	switch kind {
	case "integer":
		return "ItemInt"
	case "number":
		return "ItemFloat"
	case "boolean":
		return "ItemBool"
	default:
		return "ItemString"
	}
}

// schemaGoType maps one json-schema scalar type onto the Go type the generated
// struct field carries. An undeclared type reads as text.
func schemaGoType(kind string) string {
	switch kind {
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	default:
		return "string"
	}
}
