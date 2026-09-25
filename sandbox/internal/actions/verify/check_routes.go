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
const routesDir = utils.RoutesDir

// legacyRoutesDir is where routes lived before routeslist. A project still
// carrying it has routes no build collects any more.
const legacyRoutesDir = "sandbox/internal/routes"

// routeHandlerName is the one exported declaration a route's hand-written half
// carries, the server layer's CommandHandler.
const routeHandlerName = "InternalPureHandler"

// routeHandlerFile is the file it is declared in.
const routeHandlerFile = "InternalPureHandler.go"

// routeHandlerParams is the canonical InternalPureHandler signature the
// generated new.go closes over: the sandbox, the bound route, the route's own
// Entries and the response being written.
var routeHandlerParams = []string{"*api.Sandbox", "*api.Route", "*Entries", "*serverdeps.Response"}

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
// convention: the four files of a route package, a parsable declaration
// carrying its required keys, paths and parameters that can be matched and
// bound, Entries ids that do not collide, and a hand-written handler with the
// one signature the dispatch calls.
//
// A project with no sandbox/internal/routeslist has no server layer and nothing
// to check.
func CheckRoutes(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var violations []string

	if io.IsDir(legacyRoutesDir) {
		violations = append(violations, legacyRoutesDir+" is the old home of the routes; move each one to "+
			routesDir+"/<name>/ with a route.yaml in the methods/paths/parameters shape and an "+routeHandlerFile)
	}

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
		for _, method := range conf.Methods {
			key := method + " " + conf.Pattern() +
				" at priority " + sandbox.Deps.Stringsdeps.FormatInt(int64(conf.Priority), 10)
			if other, taken := patterns[key]; taken {
				violations = append(violations, routeViolation(name,
					"declares "+key+", which "+routesDir+"/"+other+" already declares"))
			} else {
				patterns[key] = name
			}
		}
	}

	return violations
}

// checkRouteFiles reports a route package missing any of the four files every
// route has: the declaration, the two generated files and the hand-written
// handler.
func checkRouteFiles(sandbox *api.Sandbox, io *smartio.SmartIO, name string) []string {
	var violations []string

	for _, file := range []string{"route.yaml", "new.go", "entries.go", routeHandlerFile} {
		if !io.IsFile(routesDir + "/" + name + "/" + file) {
			violations = append(violations, routeViolation(name, "has no "+file))
		}
	}

	content, err := io.ReadFile(routesDir + "/" + name + "/" + routeHandlerFile)
	if err != nil {
		return violations
	}

	parsed, err := sandbox.Deps.Goimportsdeps.Parse(string(content))
	if err != nil {
		return append(violations, routeViolation(name, routeHandlerFile+" is not parsable Go: "+err.Error()))
	}

	for _, function := range parsed.Functions {
		if function.Name != routeHandlerName {
			continue
		}
		if isRouteHandler(function) {
			return violations
		}
		return append(violations, routeViolation(name,
			routeHandlerFile+" declares "+routeHandlerName+" with another signature; the dispatch calls "+
				routeHandlerName+"(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error"))
	}

	return append(violations, routeViolation(name, routeHandlerFile+" exports no "+routeHandlerName))
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

// checkRouteDeclaration enforces the rules that survive parsing: the required
// keys, the methods, the paths, the parameters and the Entries ids they bind.
func checkRouteDeclaration(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) []string {
	var violations []string

	for _, key := range conf.Legacy {
		violations = append(violations, routeViolation(name,
			"declares `"+key+"`, which routeslist replaced: `methods` lists the methods, and `parameters` (with `fonts`) the headers and query parameters"))
	}

	if !conf.HasPriority {
		violations = append(violations, routeViolation(name, "declares no `priority`; every route declares the rung it runs on"))
	} else if conf.Priority < 0 {
		violations = append(violations, routeViolation(name,
			"declares a negative `priority`; the chain runs from zero upwards"))
	}
	if conf.ResponseType == "" {
		violations = append(violations, routeViolation(name, "declares no `response-type`; every route declares the Content-Type it answers with"))
	}

	if len(conf.Methods) == 0 {
		violations = append(violations, routeViolation(name, "declares no `methods`; a route answers one method at least"))
	}
	for _, method := range conf.Methods {
		if !contains(routeMethods, method) {
			violations = append(violations, routeViolation(name, "declares the unknown method "+method))
		}
	}

	if len(conf.Paths) == 0 {
		violations = append(violations, routeViolation(name, "declares no `paths`; a route reads one slice of the path at least"))
	}

	ids := map[string]bool{}
	for _, reserved := range utils.RouteReservedIds {
		ids[reserved] = true
	}
	claim := func(id string) {
		if !isExportedId(id) {
			violations = append(violations, routeViolation(name,
				"declares the id "+id+", which is not an exported Go name; it names a field of Entries"))
		}
		if ids[id] {
			violations = append(violations, routeViolation(name,
				"declares the id "+id+" twice, or one Entries already carries; each id names one field of Entries"))
		}
		ids[id] = true
	}

	for _, path := range conf.Paths {
		claim(path.Id)
		if path.Start < 0 {
			violations = append(violations, routeViolation(name, "declares the path "+path.Id+" with a negative `start`"))
		}
		if path.End != routeconf.LastSegment && path.End < path.Start {
			violations = append(violations, routeViolation(name,
				"declares the path "+path.Id+" with an `end` before its `start`; -1 is the last segment"))
		}
		violations = append(violations, checkRouteTrigger(sandbox, name, "path "+path.Id, path.Trigger)...)
	}

	for _, parameter := range conf.Parameters {
		claim(parameter.Id)
		if !contains(routeconf.ParameterTypes, parameter.Type) {
			violations = append(violations, routeViolation(name,
				"declares the parameter "+parameter.Id+" with the unknown type "+parameter.Type+
					" (use "+sandbox.Deps.Stringsdeps.Join(routeconf.ParameterTypes, ", ")+")"))
		}
		if len(parameter.Fonts) == 0 {
			violations = append(violations, routeViolation(name,
				"declares the parameter "+parameter.Id+" with no `fonts`; name where it is read from"))
		}
		for _, font := range parameter.Fonts {
			if !contains(routeconf.ParameterFonts, font) {
				violations = append(violations, routeViolation(name,
					"declares the parameter "+parameter.Id+" with the unknown font "+font+
						" (use "+sandbox.Deps.Stringsdeps.Join(routeconf.ParameterFonts, ", ")+")"))
			}
		}
		violations = append(violations, checkRouteTrigger(sandbox, name, "parameter "+parameter.Id, parameter.Trigger)...)
	}

	return violations
}

// checkRouteTrigger enforces what a trigger may declare: a known type, a
// value, and — for a regex — one that compiles.
func checkRouteTrigger(sandbox *api.Sandbox, name string, label string, trigger routeconf.Trigger) []string {
	if !trigger.Exists {
		return nil
	}
	if !contains(routeconf.TriggerTypes, trigger.Type) {
		return []string{routeViolation(name, "declares the "+label+" with the unknown trigger type "+trigger.Type+
			" (use "+sandbox.Deps.Stringsdeps.Join(routeconf.TriggerTypes, ", ")+")")}
	}
	if trigger.Value == "" {
		return []string{routeViolation(name, "declares the "+label+" with a trigger and no `value`")}
	}
	if trigger.Type == "regex" {
		if _, err := sandbox.Deps.Stringsdeps.MatchPattern(trigger.Value, ""); err != nil {
			return []string{routeViolation(name, "declares the "+label+" with a regex that does not compile: "+err.Error())}
		}
	}
	return nil
}

// isExportedId reports whether id reads as an exported Go name: an upper-case
// ASCII letter, then letters, digits and underscores.
func isExportedId(id string) bool {
	if id == "" || id[0] < 'A' || id[0] > 'Z' {
		return false
	}
	for _, letter := range id {
		valid := (letter >= 'a' && letter <= 'z') || (letter >= 'A' && letter <= 'Z') ||
			(letter >= '0' && letter <= '9') || letter == '_'
		if !valid {
			return false
		}
	}
	return true
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
// two contradictions a parameter entry may carry: required and defaulted, or
// required and boolean.
func checkRouteRawFields(sandbox *api.Sandbox, name string, content string) []string {
	specs, err := sandbox.Deps.Serializables.ParseYaml(content)
	if err != nil || !specs.IsObject() {
		return nil
	}

	var violations []string
	item, _ := specs.GetObjectItem("parameters")
	if item == nil || !item.IsArray() {
		return violations
	}
	size, err := item.GetArraySize()
	if err != nil {
		return violations
	}
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			continue
		}
		violations = append(violations, checkRawParameter(sandbox, name, entry)...)
	}

	return violations
}

// checkRawParameter reports the two contradictions one unparsed parameter
// entry may carry.
func checkRawParameter(sandbox *api.Sandbox, name string, entry *serializables.SerializibleObject) []string {
	if !rawBool(entry, "required") {
		return nil
	}

	label := rawString(entry, "id")

	var violations []string
	if item, _ := entry.GetObjectItem("default"); item != nil && !item.IsNull() {
		violations = append(violations, routeViolation(name,
			"declares the parameter "+label+" as both required and defaulted; the default already covers its absence"))
	}
	if rawString(entry, "type") == "boolean" || rawString(entry, "type") == "bool" {
		violations = append(violations, routeViolation(name,
			"declares the boolean parameter "+label+" as required; its absence already means false"))
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
