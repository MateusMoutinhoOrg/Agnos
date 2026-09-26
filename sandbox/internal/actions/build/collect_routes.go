package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// routesDir holds one declared route per sub-directory, the server layer's
// mirror of sandbox/internal/commands.
const routesDir = utils.RoutesDir

// CollectRoutes reads every sandbox/internal/routeslist/<name>/route.yaml and
// returns one data map per route, for the generated new.go and entries.go of
// that route and the {{range .Routes}} loop of
// sandbox/internal/generated/server/server/new.go. It is the server layer's
// CollectCommands: what the map holds is the declaration itself, which the
// dispatch reads back off Server.Routes at runtime.
//
// The list comes back ordered for running: the lowest priority first, then by
// name, so two builds of one tree always lay the chain down the same way.
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
			return nil, sandbox.Deps.Std.Errorf("routeslist/%s/route.yaml: %w", name, err)
		}

		if err := checkRouteSchema(sandbox, name, conf); err != nil {
			return nil, err
		}

		routes = append(routes, routeData(sandbox, name, conf))
	}

	sortRoutes(sandbox, routes)
	return routes, nil
}

// sortRoutes puts the routes in the order the dispatch runs them: the `before`
// phase first, then the `after` one, each by the lowest `priority` first —
// that is the whole of what the chain reads — and then by name for a stable
// tie-break. It is utils.SortRouteChain's order, over the collector's maps.
// The ordering is the collector's, so Server.Routes is already in run order
// and the dispatch only has to range.
func sortRoutes(sandbox *api.Sandbox, routes []map[string]any) {
	sandbox.Deps.Sortdeps.SliceStable(routes, func(i int, j int) bool {
		left, right := routes[i], routes[j]
		if left["After"] != right["After"] {
			return right["After"].(bool)
		}
		if left["Priority"] != right["Priority"] {
			return left["Priority"].(int) < right["Priority"].(int)
		}
		return left["Name"].(string) < right["Name"].(string)
	})
}

// checkRouteSchema refuses a json-schema that reaches outside the declared
// subset. Silently dropping $ref or oneOf would generate a validator that
// accepts what the declaration means to reject.
func checkRouteSchema(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) error {
	for _, key := range schemaUnknownKeys(conf.Body.Schema) {
		return sandbox.Deps.Std.Errorf("routeslist/%s/route.yaml: json-schema key %q is outside the supported subset", name, key)
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

// routeData is one route as the generated new.go and entries.go read it.
func routeData(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) map[string]any {
	paths := make([]map[string]any, 0, len(conf.Paths))
	for _, path := range conf.Paths {
		paths = append(paths, pathData(path))
	}
	parameters := make([]map[string]any, 0, len(conf.Parameters))
	for _, parameter := range conf.Parameters {
		parameters = append(parameters, parameterData(parameter))
	}

	return map[string]any{
		"Name":            name,
		"Methods":         conf.Methods,
		"Priority":        conf.Priority,
		"ResponseType":    conf.ResponseType,
		"Segments":        conf.Segments,
		"After":           conf.Phase == routeconf.PhaseAfter,
		"Trigger":         routeTrigger(conf),
		"Pattern":         conf.Pattern(),
		"Paths":           paths,
		"Parameters":      parameters,
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

// routeTrigger is the value of the first path trigger of a route — the one
// that names it in docs/Routes. A route whose paths declare none falls back to
// its pattern.
func routeTrigger(conf *routeconf.RouteConf) string {
	for _, path := range conf.Paths {
		if path.Trigger.Exists {
			return path.Trigger.Value
		}
	}
	return conf.Pattern()
}

// triggerConst is the api.TriggerType constant a trigger type spells.
func triggerConst(kind string) string {
	switch kind {
	case "prefix":
		return "api.PrefixTrigger"
	case "text-prefix":
		return "api.TextPrefixTrigger"
	case "suffix":
		return "api.SuffixTrigger"
	case "regex":
		return "api.RegexTrigger"
	}
	return "api.EqualTrigger"
}

// triggerData is one trigger as the generated api.Trigger literal reads it.
func triggerData(trigger routeconf.Trigger) map[string]any {
	return map[string]any{
		"Exist":      trigger.Exists,
		"Type":       triggerConst(trigger.Type),
		"Value":      trigger.Value,
		"Negate":     trigger.Negate,
		"IgnoreCase": trigger.IgnoreCase,
	}
}

// pathData is one entry of `paths` as the generated api.Path literal reads it.
func pathData(path routeconf.Path) map[string]any {
	return map[string]any{
		"Id":          path.Id,
		"Start":       path.Start,
		"End":         path.End,
		"Type":        pathTypeConst(path.Type),
		"GoType":      pathGoType(path.Type),
		"Description": path.Description,
		"Trigger":     triggerData(path.Trigger),
	}
}

// pathTypeConst is the api.PathType constant a path type spells.
func pathTypeConst(kind string) string {
	switch kind {
	case "integer":
		return "api.IntegerPath"
	case "number":
		return "api.NumberPath"
	case "uuid":
		return "api.UuidPath"
	}
	return "api.StringPath"
}

// pathGoType is the Go type of the Entries field a path binds to.
func pathGoType(kind string) string {
	switch kind {
	case "integer":
		return "int"
	case "number":
		return "float64"
	}
	return "string"
}

// parameterData is one entry of `parameters` as the generated api.Parameter
// literal and the Entries field read it.
func parameterData(parameter routeconf.Parameter) map[string]any {
	fonts := make([]string, 0, len(parameter.Fonts))
	for _, font := range parameter.Fonts {
		switch font {
		case "header":
			fonts = append(fonts, "api.HeaderParam")
		case "cookie":
			fonts = append(fonts, "api.CookieParam")
		default:
			fonts = append(fonts, "api.QueryParam")
		}
	}

	return map[string]any{
		"Id":          parameter.Id,
		"Key":         parameter.Key,
		"Fonts":       fonts,
		"Required":    parameter.Required,
		"Type":        parameterTypeConst(parameter.Type),
		"GoType":      parameterGoType(parameter.Type),
		"Default":     parameter.Default,
		"HasDefault":  parameter.HasDefault,
		"Description": parameter.Description,
		"Trigger":     triggerData(parameter.Trigger),
	}
}

// parameterTypeConst is the api.ParameterType constant a parameter type spells.
func parameterTypeConst(kind string) string {
	switch kind {
	case "integer":
		return "api.IntegerType"
	case "number":
		return "api.NumberType"
	case "boolean":
		return "api.BooleanType"
	case "datetime":
		return "api.DateTimeType"
	case "string-array":
		return "api.StringArrayType"
	case "integer-array":
		return "api.IntegerArrayType"
	}
	return "api.StringType"
}

// parameterGoType is the Go type of the Entries field a parameter binds to.
func parameterGoType(kind string) string {
	switch kind {
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "string-array":
		return "[]string"
	case "integer-array":
		return "[]int"
	}
	return "string"
}

// bodyData is the route's body declaration as the generated new.go and
// entries.go read it.
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
		"IsForm":       body.Type == "form",
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
	case "form":
		return "map[string][]string"
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
