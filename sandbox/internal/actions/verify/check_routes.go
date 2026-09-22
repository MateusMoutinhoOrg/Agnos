package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	serializables "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// routesDir holds one declared route per sub-directory.
const routesDir = "sandbox/internal/routes"

// routeHandlerName is the one exported declaration a route's hand-written half
// carries, the server layer's CommandHandler.
const routeHandlerName = "RouteHandler"

// routeHandlerParams is the canonical RouteHandler signature the generated
// new.go closes over: the sandbox, the bound route and the response being
// written.
var routeHandlerParams = []string{"*api.Sandbox", "*api.Route", "serverdeps.Response"}

// routeMethods is every http method a route may declare.
var routeMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// schemaKeys is the subset of JSON Schema a route's body may declare.
var schemaKeys = []string{
	"type", "properties", "required", "additionalProperties", "items", "enum",
	"const", "minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum",
	"minLength", "maxLength", "pattern", "minItems", "maxItems", "uniqueItems",
	"format", "nullable",
}

// CheckRoutes enforces the shape the server layer's generators read by
// convention: the three files of a route package, a parsable declaration, path
// segments that can be matched, fields that agree with each other, and a
// hand-written handler with the one signature the dispatch calls.
//
// A project with no sandbox/internal/routes has no server layer and nothing to
// check.
func CheckRoutes(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var violations []string

	if !io.IsDir(routesDir) {
		return violations
	}

	patterns := map[string]string{}

	for _, dir := range io.ListDirs(routesDir) {
		name := routeNameOf(sandbox, dir)
		if name == "" {
			continue
		}

		violations = append(violations, checkRouteFiles(sandbox, io, name)...)

		content, err := io.ReadFile(routesDir + "/" + name + "/route.yaml")
		if err != nil {
			continue
		}

		conf, err := routeconf.New(sandbox, string(content))
		if err != nil {
			violations = append(violations, routeViolation(name, "route.yaml does not parse: "+err.Error()))
			continue
		}

		violations = append(violations, checkRouteDeclaration(sandbox, name, conf)...)
		violations = append(violations, checkRouteRawFields(sandbox, name, string(content))...)
		violations = append(violations, checkRouteSchemaKeys(sandbox, name, conf)...)

		// Two routes may share a method and a pattern as long as they sit
		// on different rungs of the chain — that is what a middleware in
		// front of a route is. Sharing a rung as well is the ambiguity: the
		// two would run in an order nothing declares.
		key := conf.Method + " " + conf.Pattern() +
			" at priority " + sandbox.Deps.Stringsdeps.FormatInt(int64(conf.Priority), 10)
		if other, taken := patterns[key]; taken {
			violations = append(violations, routeViolation(name,
				"declares "+key+", which routes/"+other+" already declares"))
		} else {
			patterns[key] = name
		}
	}

	return violations
}

// checkRouteFiles reports a route package missing any of the three files every
// route has: the declaration, the generated constructor and the hand-written
// handler.
func checkRouteFiles(sandbox *api.Sandbox, io *smartio.SmartIO, name string) []string {
	var violations []string

	for _, file := range []string{"route.yaml", "new.go", "handler.go"} {
		if !io.IsFile(routesDir + "/" + name + "/" + file) {
			violations = append(violations, routeViolation(name, "has no "+file))
		}
	}

	content, err := io.ReadFile(routesDir + "/" + name + "/handler.go")
	if err != nil {
		return violations
	}

	parsed, err := sandbox.Deps.Goimportsdeps.Parse(string(content))
	if err != nil {
		return append(violations, routeViolation(name, "handler.go is not parsable Go: "+err.Error()))
	}

	for _, function := range parsed.Functions {
		if function.Name != routeHandlerName {
			continue
		}
		if isRouteHandler(function) {
			return violations
		}
		return append(violations, routeViolation(name,
			"handler.go declares "+routeHandlerName+" with another signature; the dispatch calls "+
				routeHandlerName+"(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error"))
	}

	return append(violations, routeViolation(name, "handler.go exports no "+routeHandlerName))
}

// isRouteHandler reports whether one parsed declaration is the route handler:
// a plain exported func of the canonical name, parameters and one error
// result. The status is what the handler writes on the response, never what it
// returns — what it returns is the failure it did not answer itself.
func isRouteHandler(function goimportsdeps.Function) bool {
	if function.Receiver != "" || len(function.Params) != len(routeHandlerParams) {
		return false
	}
	for i, kind := range routeHandlerParams {
		if function.Params[i].Type != kind {
			return false
		}
	}
	return len(function.Results) == 1 && function.Results[0].Type == "error"
}

// checkRouteDeclaration enforces the rules that survive parsing: the method,
// the segments and the names a route binds.
func checkRouteDeclaration(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) []string {
	var violations []string

	known := false
	for _, method := range routeMethods {
		if method == conf.Method {
			known = true
			break
		}
	}
	if !known {
		violations = append(violations, routeViolation(name, "declares the unknown method "+conf.Method))
	}

	if len(conf.Paths) == 0 {
		violations = append(violations, routeViolation(name, "declares no `paths`; a route needs at least one segment"))
	}

	if conf.Priority < 0 {
		violations = append(violations, routeViolation(name,
			"declares a negative `priority`; the chain runs from zero upwards"))
	}

	triggers := 0
	for i, segment := range conf.Paths {
		if segment.Field == nil {
			triggers++
			violations = append(violations, checkRouteTrigger(sandbox, name, segment, i == len(conf.Paths)-1)...)
			continue
		}
		violations = append(violations, checkRouteCapture(name, *segment.Field, i == len(conf.Paths)-1)...)
	}
	if len(conf.Paths) > 0 && triggers == 0 {
		violations = append(violations, routeViolation(name,
			"declares no `identifier` segment; a route is named by the first literal segment of its path"))
	}

	violations = append(violations, checkRoutePrefix(sandbox, name, conf)...)
	violations = append(violations, checkRouteNames(sandbox, name, conf)...)
	violations = append(violations, checkRouteRestName(sandbox, name, conf)...)
	violations = append(violations, checkRouteFieldIdentifiers(sandbox, name, conf)...)

	return violations
}

// checkRoutePrefix keeps the two ways of taking the rest of a path apart: a
// prefix trigger leaves every segment after it unmatched, and so does an array
// capture, so a route declaring both says twice what it binds once.
func checkRoutePrefix(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) []string {
	prefix := false
	for _, segment := range conf.Paths {
		if segment.Field == nil && segment.StartsWith {
			prefix = true
		}
	}
	if !prefix || utils.RouteRestIndex(conf.Paths) < 0 {
		return nil
	}
	return []string{routeViolation(name,
		"declares both a `starts-with-identifier` and an array capture; each one already takes every segment left in the path")}
}

// checkRouteFieldIdentifiers enforces the two value conditions a header or a
// query parameter may carry. They are what puts the field in the match, so a
// field declaring both says two things about one value.
func checkRouteFieldIdentifiers(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) []string {
	var violations []string

	origins := []struct {
		label  string
		fields []routeconf.Field
	}{
		{"header", conf.Headers},
		{"query", conf.Params},
	}

	for _, origin := range origins {
		for _, field := range origin.fields {
			if field.Identifier != "" && field.StartsWith != "" {
				violations = append(violations, routeViolation(name,
					"declares the "+origin.label+" field "+field.Key+
						" with both an `identifier` and a `starts-with-identifier`; a value matches one condition or the other"))
			}
		}
	}

	return violations
}

// checkRouteRestName keeps the segment taking the rest of the path out of the
// name merging every other origin takes part in: what it binds is a list of
// values, which a header or a query parameter of the same name could not go on
// filling.
func checkRouteRestName(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) []string {
	at := utils.RouteRestIndex(conf.Paths)
	if at < 0 {
		return nil
	}

	var violations []string
	key := conf.Paths[at].Field.Key
	origins := []struct {
		label  string
		fields []routeconf.Field
	}{
		{"header", conf.Headers},
		{"query", conf.Params},
	}

	for _, origin := range origins {
		if utils.FindRouteField(sandbox, origin.fields, key) >= 0 {
			violations = append(violations, routeViolation(name,
				"declares "+key+" as both the segment taking the rest of the path and a "+origin.label+" field; the two cannot fill one field"))
		}
	}

	return violations
}

// checkRouteTrigger enforces what a literal segment may spell. Both kinds
// carry the leading slash. An `identifier` is exactly one segment; a
// `starts-with-identifier` may spell several, because it names the head of a
// whole subtree rather than one place in it — and it is the last entry of
// `paths`, since everything after it is unmatched by definition.
func checkRouteTrigger(sandbox *api.Sandbox, name string, segment routeconf.Segment, is_last bool) []string {
	label := "path identifier"
	if segment.StartsWith {
		label = "path starts-with-identifier"
	}

	if segment.Identifier == "" {
		return []string{routeViolation(name, "declares an empty "+label)}
	}
	if !sandbox.Deps.Stringsdeps.HasPrefix(segment.Identifier, "/") {
		return []string{routeViolation(name, "declares the "+label+" "+segment.Identifier+
			", which does not start with '/'")}
	}

	if segment.StartsWith {
		if !is_last {
			return []string{routeViolation(name, "declares the "+label+" "+segment.Identifier+
				" before the end of `paths`; a starts-with-identifier leaves every segment after it unmatched, so nothing may follow it")}
		}
		if segment.Identifier != "/" && sandbox.Deps.Stringsdeps.HasSuffix(segment.Identifier, "/") {
			return []string{routeViolation(name, "declares the "+label+" "+segment.Identifier+
				", which holds a trailing slash; spell the segments it fixes and no more")}
		}
		return nil
	}

	if segment.Identifier == "/" {
		return nil
	}
	if sandbox.Deps.Stringsdeps.Contains(sandbox.Deps.Stringsdeps.TrimPrefix(segment.Identifier, "/"), "/") {
		return []string{routeViolation(name, "declares the "+label+" "+segment.Identifier+
			", which holds an inner or trailing slash; an identifier is one segment")}
	}
	return nil
}

// checkRouteCapture enforces what a captured segment may declare: it is always
// present when the route matches, so it is always required and never defaults.
// `array: true` makes it take every segment left in the path, which only reads
// as a suffix on the last entry of `paths`.
func checkRouteCapture(name string, field routeconf.Field, is_last bool) []string {
	var violations []string

	if !field.Required {
		violations = append(violations, routeViolation(name,
			"declares the captured segment "+field.Key+" without `required: true`; a captured segment is always present"))
	}
	if field.Array && !is_last {
		violations = append(violations, routeViolation(name,
			"declares the captured segment "+field.Key+" as an array before the end of `paths`; only the last segment takes the rest of the path"))
	}
	if field.HasDefault {
		violations = append(violations, routeViolation(name,
			"declares a default for the captured segment "+field.Key+", which is always present"))
	}

	return violations
}

// checkRouteNames enforces that a name is declared once per origin, and that
// the origins declaring the same name agree on its type — one name binds one
// field of Route.Items, so two origins cannot disagree about what it holds.
func checkRouteNames(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) []string {
	var violations []string
	types := map[string]string{}

	origins := []struct {
		label  string
		fields []routeconf.Field
	}{
		{"path", captureFields(conf)},
		{"header", conf.Headers},
		{"query", conf.Params},
	}

	for _, origin := range origins {
		seen := map[string]bool{}
		for _, field := range origin.fields {
			if seen[field.Key] {
				violations = append(violations, routeViolation(name,
					"declares the "+origin.label+" field "+field.Key+" twice"))
				continue
			}
			seen[field.Key] = true

			if kind, known := types[field.Key]; known && kind != field.Type {
				violations = append(violations, routeViolation(name,
					"declares "+field.Key+" as "+kind+" in one origin and as "+field.Type+" in another"))
				continue
			}
			types[field.Key] = field.Type
		}
	}

	return violations
}

// captureFields is every captured segment of a route, in declaration order.
func captureFields(conf *routeconf.RouteConf) []routeconf.Field {
	var fields []routeconf.Field
	for _, segment := range conf.Paths {
		if segment.Field != nil {
			fields = append(fields, *segment.Field)
		}
	}
	return fields
}

// checkRouteSchemaKeys reports a json-schema declared on a body that carries
// none, and every keyword outside the supported subset, whatever depth it sits
// at.
func checkRouteSchemaKeys(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) []string {
	var violations []string

	if conf.Body.HasSchema && conf.Body.Type != "json" {
		violations = append(violations, routeViolation(name,
			"declares a json-schema on a "+conf.Body.Type+" body; only `type: json` carries one"))
	}

	for _, key := range schemaUnknown(conf.Body.Schema) {
		violations = append(violations, routeViolation(name,
			"declares the json-schema key "+key+", which is outside the supported subset ("+
				sandbox.Deps.Stringsdeps.Join(schemaKeys, ", ")+")"))
	}

	return violations
}

// schemaUnknown walks a schema tree and gathers every key outside the subset.
func schemaUnknown(schema *routeconf.Schema) []string {
	if schema == nil {
		return nil
	}
	unknown := schema.Unknown
	unknown = append(unknown, schemaUnknown(schema.Items)...)
	for _, property := range schema.Properties {
		unknown = append(unknown, schemaUnknown(property.Schema)...)
	}
	return unknown
}

// checkRouteRawFields reads the declaration a second time, unparsed, for the
// two rules the parser resolves away: it drops `required` from a boolean and
// from a field that carries a default, so by the time a RouteConf exists the
// contradiction is gone. The file is what has to be right.
func checkRouteRawFields(sandbox *api.Sandbox, name string, content string) []string {
	specs, err := sandbox.Deps.Serializables.ParseYaml(content)
	if err != nil || !specs.IsObject() {
		return nil
	}

	var violations []string
	for _, key := range []string{"paths", "headers", "params"} {
		item, _ := specs.GetObjectItem(key)
		if item == nil || !item.IsArray() {
			continue
		}
		size, err := item.GetArraySize()
		if err != nil {
			continue
		}
		for i := 0; i < size; i++ {
			entry := item.GetArrayItem(i)
			if entry == nil || !entry.IsObject() {
				continue
			}
			violations = append(violations, checkRawField(sandbox, name, key, entry)...)
		}
	}

	return violations
}

// checkRawField reports the two contradictions one unparsed field entry may
// carry.
func checkRawField(sandbox *api.Sandbox, name string, origin string, entry *serializables.SerializibleObject) []string {
	if !rawBool(entry, "required") {
		return nil
	}

	label := rawString(entry, "name")
	if label == "" {
		label = rawString(entry, "identifier")
	}
	if label == "" {
		label = rawString(entry, "starts-with-identifier")
	}

	var violations []string
	if item, _ := entry.GetObjectItem("default"); item != nil && !item.IsNull() {
		violations = append(violations, routeViolation(name,
			"declares the "+origin+" field "+label+" as both required and defaulted; the default already covers its absence"))
	}
	if rawString(entry, "type") == "boolean" || rawString(entry, "type") == "bool" {
		violations = append(violations, routeViolation(name,
			"declares the boolean "+origin+" field "+label+" as required; its absence already means false"))
	}
	return violations
}

func rawString(entry *serializables.SerializibleObject, key string) string {
	item, _ := entry.GetObjectItem(key)
	if item == nil || !item.IsString() {
		return ""
	}
	value, err := item.GetString()
	if err != nil {
		return ""
	}
	return value
}

func rawBool(entry *serializables.SerializibleObject, key string) bool {
	item, _ := entry.GetObjectItem(key)
	if item == nil || !item.IsBool() {
		return false
	}
	value, err := item.GetBool()
	if err != nil {
		return false
	}
	return value
}

// routeNameOf is the last segment of a listed route directory.
func routeNameOf(sandbox *api.Sandbox, path string) string {
	parts := sandbox.Deps.Stringsdeps.Split(path, "/")
	return parts[len(parts)-1]
}

// routeViolation words one violation the same way for every rule.
func routeViolation(name string, reason string) string {
	return routesDir + "/" + name + " " + reason
}
