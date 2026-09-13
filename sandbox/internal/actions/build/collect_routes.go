package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// routesDir holds one declared route per sub-directory, the server layer's
// mirror of sandbox/internal/commands.
const routesDir = "sandbox/internal/routes"

// CollectRoutes reads every sandbox/internal/routes/<name>/route.yaml and
// returns one data map per route, for the generated
// sandbox/internal/routes/<name>/new.go and the {{range .Routes}} loop of
// sandbox/binds/server.go. It is the server layer's CollectCommands: what the
// map holds is the declaration itself, which the dispatch reads back off
// sandbox.Routes at runtime, so nothing here is a Go spelling of the routing.
//
// The list comes back ordered for matching, not in listing order: the most
// specific route first, so a route fixing "/" can never swallow "/home".
func CollectRoutes(sandbox *api.Sandbox, io *smartio.SmartIO) ([]map[string]any, error) {
	var routes []map[string]any

	for _, dir := range io.ListDirs(routesDir) {
		name := lastSegmentOf(sandbox, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile(routesDir + "/" + name + "/route.yaml")
		if err != nil {
			continue
		}

		conf, err := routeconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("routes/%s/route.yaml: %w", name, err)
		}

		if err := checkRouteSchema(sandbox, name, conf); err != nil {
			return nil, err
		}

		if err := checkRoutePaths(sandbox, name, conf); err != nil {
			return nil, err
		}

		routes = append(routes, routeData(sandbox, name, conf))
	}

	sortRoutes(sandbox, routes)
	return routes, nil
}

// sortRoutes puts the routes in the order the dispatch tests them: the route
// fixing the most literal segments first, then the one whose literals spell
// the most characters, then the one of fixed length before the one taking the
// rest of the path, then by pattern for a stable tie-break. The ordering is
// the collector's, so sandbox.Routes is already in match order and the
// dispatch only has to range.
func sortRoutes(sandbox *api.Sandbox, routes []map[string]any) {
	sandbox.Deps.Sortdeps.SliceStable(routes, func(i int, j int) bool {
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
func checkRouteSchema(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) error {
	for _, key := range schemaUnknownKeys(conf.Body.Schema) {
		return sandbox.Deps.Std.Errorf("routes/%s/route.yaml: json-schema key %q is outside the supported subset", name, key)
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

// checkRoutePaths refuses a `paths` the dispatch cannot match: the segment
// taking the rest of the path only reads as a suffix if it is the last one,
// and what it binds is a list, which no header or query parameter of the same
// name could go on filling.
func checkRoutePaths(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) error {
	at := utils.RouteRestIndex(conf.Paths)
	if at < 0 {
		return nil
	}

	if at != len(conf.Paths)-1 {
		return sandbox.Deps.Std.Errorf("routes/%s/route.yaml: the captured segment %q is an array before the end of `paths`; only the last segment takes the rest of the path", name, conf.Paths[at].Field.Key)
	}

	key := conf.Paths[at].Field.Key
	for _, origin := range []struct {
		label  string
		fields []routeconf.Field
	}{{"header", conf.Headers}, {"query parameter", conf.Params}} {
		if utils.FindRouteField(sandbox, origin.fields, key) >= 0 {
			return sandbox.Deps.Std.Errorf("routes/%s/route.yaml: %q is both the segment taking the rest of the path and a declared %s; the two cannot fill one field", name, key, origin.label)
		}
	}

	return nil
}

// routeData is one route as the generated new.go reads it, plus the three keys
// sortRoutes orders the set by.
func routeData(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) map[string]any {
	paths := make([]map[string]any, 0, len(conf.Paths))
	for _, segment := range conf.Paths {
		paths = append(paths, segmentData(sandbox, segment))
	}
	headers := make([]map[string]any, 0, len(conf.Headers))
	for _, field := range conf.Headers {
		headers = append(headers, routeFieldData(sandbox, field))
	}
	params := make([]map[string]any, 0, len(conf.Params))
	for _, field := range conf.Params {
		params = append(params, routeFieldData(sandbox, field))
	}

	return map[string]any{
		"Name":            name,
		"Method":          conf.Method,
		"Trigger":         routeTrigger(conf),
		"Pattern":         conf.Pattern(),
		"IdentifierCount": conf.IdentifierCount(),
		"IdentifierLen":   conf.IdentifierLen(),
		"HasRest":         utils.RouteRestIndex(conf.Paths) >= 0,
		"Paths":           paths,
		"Headers":         headers,
		"Params":          params,
		"Category":        conf.Category,
		"Help":            conf.Help,
		"LongDescription": conf.LongDescription,
		"Examples":        conf.Examples,
		"Hidden":          conf.Hidden,
		"Body":            bodyData(sandbox, conf),
		"HasBody":         conf.Body.Type != routeconf.BodyNone,
		"SchemaJson":      conf.SchemaJson(),
		"BodyStructs":     bodyStructs(sandbox, conf),
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

// segmentData is one entry of `paths` as the generated api.RoutePath literal
// reads it: a literal segment carries its identifier, a capture its field.
func segmentData(sandbox *api.Sandbox, segment routeconf.Segment) map[string]any {
	if segment.Field == nil {
		return map[string]any{
			"HasField":   false,
			"Identifier": segment.Identifier,
		}
	}
	return map[string]any{
		"HasField": true,
		"Field":    routeFieldData(sandbox, *segment.Field),
	}
}

// routeFieldData is one declared field as the generated api.RouteField literal
// reads it — the server layer's fieldData, with the header name or the query
// key standing where a flag's identifiers would.
func routeFieldData(sandbox *api.Sandbox, field routeconf.Field) map[string]any {
	return map[string]any{
		"Key":         field.Key,
		"Type":        field.Type,
		"IsArray":     field.Array,
		"Required":    field.Required,
		"HasDefault":  field.HasDefault,
		"Default":     field.Default,
		"Description": field.Description,
		"Examples":    field.Examples,
		"Min":         routeNumberLabel(sandbox, field.Type, field.Min, field.HasMin),
		"Max":         routeNumberLabel(sandbox, field.Type, field.Max, field.HasMax),
		"HasMin":      field.HasMin,
		"HasMax":      field.HasMax,
	}
}

// routeNumberLabel renders a min/max bound as the literal it has in
// route.yaml ("" when the bound is unset).
func routeNumberLabel(sandbox *api.Sandbox, kind string, value float64, has bool) string {
	if !has {
		return ""
	}
	if kind == "int" {
		return sandbox.Deps.Stringsdeps.FormatInt(int64(value), 10)
	}
	return sandbox.Deps.Stringsdeps.FormatFloat(value, 'g', -1, 64)
}

// bodyData is the route's body declaration as the generated new.go reads it.
func bodyData(sandbox *api.Sandbox, conf *routeconf.RouteConf) map[string]any {
	body := conf.Body
	return map[string]any{
		"Type":         body.Type,
		"Required":     body.Required,
		"MaxBytes":     body.MaxBytes,
		"MaxBytesText": sandbox.Deps.Stringsdeps.FormatInt(int64(body.MaxBytes), 10),
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
// generated new.go declares: Body for the root object, Body<Path> for a
// nested object and Body<Path>Item for the object an array holds.
func bodyStructs(sandbox *api.Sandbox, conf *routeconf.RouteConf) []map[string]any {
	body := conf.Body
	if body.Type != "json" || !body.HasSchema || body.Schema == nil || body.Schema.Type != "object" {
		return nil
	}
	return appendBodyStruct(sandbox, nil, "Body", body.Schema)
}

// appendBodyStruct emits one struct for schema and, depth first, one for every
// object nested under it. The parent is appended before its children so the
// generated file reads outside in.
func appendBodyStruct(sandbox *api.Sandbox, structs []map[string]any, name string, schema *routeconf.Schema) []map[string]any {
	fields := make([]map[string]any, 0, len(schema.Properties))
	var nested []*routeconf.Schema
	var nested_names []string

	for _, property := range schema.Properties {
		go_field := exportedName(sandbox, property.Name)
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
		structs = appendBodyStruct(sandbox, structs, nested_names[i], child)
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
