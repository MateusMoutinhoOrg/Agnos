package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	serializabledeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializabledeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// routesDir is the tree the routes are declared in, at any depth.
const routesDir = utils.RoutesDir

// routeHandlerName is the one exported declaration a route's hand-written half
// carries, the server layer's twin of a command's Handle.
const routeHandlerName = "Handle"

// routeHandlerFile is the file it is declared in.
const routeHandlerFile = "handler.go"

// routeHandlerParams is the canonical Handle signature the
// generated.new.go closes over: the sandbox, the request's shared RouteProps,
// the route's own Input and the response being written.
var routeHandlerParams = []string{"*api.Sandbox", "*routeprops.RouteProps", "*Input", "*serverdeps.Response"}

// legacyRoutePropsParam is the props parameter a handler took while
// RouteProps was declared in sandbox/api.
const legacyRoutePropsParam = "*api.RouteProps"

// routeMethods is every http method a route may declare beside ANY, which
// stands alone.
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
// bound, Input ids that do not collide, and a hand-written handler with the
// one signature the dispatch calls.
//
// A project with no sandbox/internal/routes has no server layer and nothing
// to check.
func CheckRoutes(sandbox *api.Sandbox, io *stagedfs.StagedFS) []string {
	var violations []string

	if !io.IsDir(routesDir) {
		return violations
	}

	patterns := map[string]string{}

	violations = append(violations, CheckUnitTree(sandbox, io, routesDir, utils.RouteConfFile, "route",
		utils.GeneratedNames(sandbox, utils.UnitNewFile, utils.UnitInputFile, routeHandlerFile))...)

	for _, unit := range utils.RouteDirs(sandbox, io) {
		dir := unit.Dir

		violations = append(violations, checkRouteFiles(sandbox, io, dir)...)

		content, err := io.ReadFile(dir + "/route.yaml")
		if err != nil {
			continue
		}

		conf, err := routeconf.New(sandbox, string(content))
		if err != nil {
			violations = append(violations, routeViolation(dir, "route.yaml does not parse: "+err.Error()))
			continue
		}

		violations = append(violations, checkRouteDeclaration(sandbox, dir, conf)...)
		violations = append(violations, checkRouteRawFields(sandbox, dir, string(content))...)
		violations = append(violations, checkRouteSchemaKeys(sandbox, dir, conf)...)

		// Two routes may share a method and a pattern as long as they sit
		// on different rungs of the chain — that is what a middleware in
		// front of a route is. Sharing a rung as well is the ambiguity: the
		// two would run in an order nothing declares.
		for _, method := range conf.Methods {
			key := method + " " + conf.Pattern() +
				" at priority " + sandbox.Deps.StringsDeps.FormatInt(int64(conf.Priority), 10)
			if other, taken := patterns[key]; taken {
				violations = append(violations, routeViolation(dir,
					"declares "+key+", which "+other+" already declares"))
			} else {
				patterns[key] = dir
			}
		}
	}

	return violations
}

// checkRouteFiles reports a route package missing any of the four files every
// route has: the declaration, the two generated files (generated.new.go and
// generated.input.go, or their old names on a tree no build has moved yet)
// and the hand-written handler (generated.handler.go for a route a group
// renders whole).
func checkRouteFiles(sandbox *api.Sandbox, io *stagedfs.StagedFS, dir string) []string {
	var violations []string

	if !io.IsFile(dir + "/route.yaml") {
		violations = append(violations, routeViolation(dir, "has no route.yaml"))
	}
	for _, file := range []string{utils.UnitNewFile, utils.UnitInputFile} {
		if utils.GeneratedPath(sandbox, io, dir, file) == "" {
			violations = append(violations, routeViolation(dir, "has no "+utils.GeneratedFile(sandbox, file)))
		}
	}

	handler := utils.GeneratedPath(sandbox, io, dir, routeHandlerFile)
	if handler == "" {
		return append(violations, routeViolation(dir, "has no "+routeHandlerFile))
	}
	handlerFile := lastSegment(sandbox, handler)
	content, err := io.ReadFile(handler)
	if err != nil {
		return violations
	}

	parsed, err := sandbox.Deps.GoimportsDeps.Parse(string(content))
	if err != nil {
		return append(violations, routeViolation(dir, handlerFile+" is not parsable Go: "+err.Error()))
	}

	for _, function := range parsed.Functions {
		if function.Name != routeHandlerName {
			continue
		}
		if isRouteHandler(function) {
			return violations
		}
		return append(violations, routeViolation(dir,
			handlerFile+" declares "+routeHandlerName+" with another signature; the dispatch calls "+
				routeHandlerName+"(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error"+
				legacyPropsHint(function, legacyRoutePropsParam, utils.RoutePropsDir)))
	}

	return append(violations, routeViolation(dir, handlerFile+" exports no "+routeHandlerName))
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

// legacyPropsHint is what a handler still taking its props from sandbox/api is
// told on top of the signature: where the struct lives now. "" for any other
// mismatch.
func legacyPropsHint(function goimportsdeps.Function, legacy string, dir string) string {
	if len(function.Params) < 2 || function.Params[1].Type != legacy {
		return ""
	}
	return " (the props struct moved out of sandbox/api to " + dir + ": take it from there and import that package)"
}

// checkRouteDeclaration enforces the rules that survive parsing: the required
// keys, the methods, the paths, the parameters and the Input ids they bind.
func checkRouteDeclaration(sandbox *api.Sandbox, dir string, conf *routeconf.RouteConf) []string {
	var violations []string

	for _, key := range conf.Legacy {
		if key == "phase" {
			violations = append(violations, routeViolation(dir,
				"declares `phase`, which is gone: every route is a rung of the one chain, so drop the key"))
			continue
		}
		violations = append(violations, routeViolation(dir,
			"declares `"+key+"`, which the current route.yaml shape replaced: `methods` lists the methods, and `parameters` (with `sources`) the headers and query parameters"))
	}

	if !conf.HasPriority {
		violations = append(violations, routeViolation(dir, "declares no `priority`; every route declares the rung it runs on"))
	} else if conf.Priority < 0 {
		violations = append(violations, routeViolation(dir,
			"declares a negative `priority`; the chain runs from zero upwards"))
	}
	if conf.ResponseType == "" {
		violations = append(violations, routeViolation(dir, "declares no `response-type`; every route declares the Content-Type it answers with"))
	} else if err := utils.ValidateMediaType(sandbox, conf.ResponseType); err != nil {
		violations = append(violations, routeViolation(dir, err.Error()))
	}

	if len(conf.Methods) == 0 {
		violations = append(violations, routeViolation(dir, "declares no `methods`; a route answers one method at least"))
	}
	for _, method := range conf.Methods {
		if method == routeconf.AnyMethod {
			if len(conf.Methods) > 1 {
				violations = append(violations, routeViolation(dir,
					"declares "+routeconf.AnyMethod+" beside other methods; "+routeconf.AnyMethod+" already accepts every one"))
			}
			continue
		}
		if !contains(routeMethods, method) {
			violations = append(violations, routeViolation(dir, "declares the unknown method "+method))
		}
	}

	if conf.HasSegments && conf.Segments < 1 {
		violations = append(violations, routeViolation(dir,
			"declares `segments` below 1; leave it out for a route that takes any count"))
	}
	if !contains(routeconf.BodyTypes, conf.Body.Type) {
		violations = append(violations, routeViolation(dir, "declares the unknown body type "+conf.Body.Type))
	}

	if len(conf.Paths) == 0 {
		violations = append(violations, routeViolation(dir, "declares no `paths`; a route reads one slice of the path at least"))
	}

	ids := map[string]bool{}
	for _, reserved := range utils.RouteReservedIds {
		ids[reserved] = true
	}
	claim := func(id string) {
		if !isExportedId(id) {
			violations = append(violations, routeViolation(dir,
				"declares the id "+id+", which is not an exported ASCII Go name (an uppercase ASCII letter, then ASCII letters and digits); it names a field of Input"))
		}
		if ids[id] {
			violations = append(violations, routeViolation(dir,
				"declares the id "+id+" twice, or one Input already carries; each id names one field of Input"))
		}
		ids[id] = true
	}

	for _, path := range conf.Paths {
		claim(path.Id)
		if path.Start < 0 {
			violations = append(violations, routeViolation(dir, "declares the path "+path.Id+" with a negative `start`"))
		}
		if path.End != routeconf.LastSegment && path.End < path.Start {
			violations = append(violations, routeViolation(dir,
				"declares the path "+path.Id+" with an `end` before its `start`; -1 is the last segment"))
		}
		if !contains(routeconf.PathTypes, path.Type) {
			violations = append(violations, routeViolation(dir,
				"declares the path "+path.Id+" with the unknown type "+path.Type+
					" (use "+sandbox.Deps.StringsDeps.Join(routeconf.PathTypes, ", ")+")"))
		} else if path.Type != routeconf.DefaultPathType && path.Start != path.End {
			violations = append(violations, routeViolation(dir,
				"declares the path "+path.Id+" of type "+path.Type+" over more than one segment; a typed path reads one, so `start` and `end` are the same index"))
		}
		violations = append(violations, checkRouteTrigger(sandbox, dir, "path "+path.Id, path.Trigger)...)
	}

	for _, parameter := range conf.Parameters {
		claim(parameter.Id)
		if !contains(routeconf.ParameterTypes, parameter.Type) {
			violations = append(violations, routeViolation(dir,
				"declares the parameter "+parameter.Id+" with the unknown type "+parameter.Type+
					" (use "+sandbox.Deps.StringsDeps.Join(routeconf.ParameterTypes, ", ")+")"))
		}
		if len(parameter.Sources) == 0 {
			violations = append(violations, routeViolation(dir,
				"declares the parameter "+parameter.Id+" with no `sources`; name where it is read from"))
		}
		for _, source := range parameter.Sources {
			if !contains(routeconf.ParameterSources, source) {
				violations = append(violations, routeViolation(dir,
					"declares the parameter "+parameter.Id+" with the unknown source "+source+
						" (use "+sandbox.Deps.StringsDeps.Join(routeconf.ParameterSources, ", ")+")"))
			}
		}
		violations = append(violations, checkRouteTrigger(sandbox, dir, "parameter "+parameter.Id, parameter.Trigger)...)
	}

	return violations
}

// checkRouteTrigger enforces what a trigger may declare on a route; see
// triggerProblem.
func checkRouteTrigger(sandbox *api.Sandbox, dir string, label string, trigger triggerconf.Trigger) []string {
	if problem := triggerProblem(sandbox, trigger, false); problem != "" {
		return []string{routeViolation(dir, "declares the "+label+" with "+problem)}
	}
	return nil
}

// triggerProblem is what is wrong with one declared trigger, "" when nothing
// is: a known type, a value — the `values` of a one-of — and, for a regex, one
// that compiles. A prefix with no value matches everything, which is what a
// command middleware declares, so emptyPrefix lets it through.
func triggerProblem(sandbox *api.Sandbox, trigger triggerconf.Trigger, emptyPrefix bool) string {
	if !trigger.Set {
		return ""
	}
	if !contains(triggerconf.TriggerTypes, trigger.Type) {
		return "the unknown trigger type " + trigger.Type +
			" (use " + sandbox.Deps.StringsDeps.Join(triggerconf.TriggerTypes, ", ") + ")"
	}
	if trigger.Type == triggerconf.OneOf {
		if len(trigger.Values) == 0 {
			return "a one-of trigger and no `values`"
		}
		return ""
	}
	if trigger.Value == "" && !(emptyPrefix && trigger.Type == "prefix") {
		return "a trigger and no `value`"
	}
	if trigger.Type == triggerconf.Regex {
		if _, err := sandbox.Deps.StringsDeps.MatchPattern(trigger.Value, ""); err != nil {
			return "a regex that does not compile: " + err.Error()
		}
	}
	return ""
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

// checkRouteSchemaKeys reports a schema declared on a body that carries none,
// or under the key of the other type (`json-schema` on a form body), a form
// schema that is not flat, and every keyword outside the supported subset,
// whatever depth it sits at.
func checkRouteSchemaKeys(sandbox *api.Sandbox, dir string, conf *routeconf.RouteConf) []string {
	var violations []string

	expected := routeconf.SchemaKeyOf(conf.Body.Type)
	for _, key := range conf.Body.SchemaKeys {
		switch {
		case expected == "":
			violations = append(violations, routeViolation(dir,
				"declares a "+key+" on a "+conf.Body.Type+" body; only `type: json` (json-schema) and `type: form` (form-schema) carry one"))
		case key != expected:
			violations = append(violations, routeViolation(dir,
				"declares a "+key+" on a "+conf.Body.Type+" body, which carries its schema under "+expected))
		}
	}

	if conf.Body.Type == "form" {
		for _, violation := range utils.FormSchemaViolations(conf.Body.Schema) {
			violations = append(violations, routeViolation(dir, "declares a form-schema that is not flat: "+violation))
		}
	}

	for _, key := range schemaUnknown(conf.Body.Schema) {
		violations = append(violations, routeViolation(dir,
			"declares the schema key "+key+", which is outside the supported subset ("+
				sandbox.Deps.StringsDeps.Join(schemaKeys, ", ")+")"))
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
func checkRouteRawFields(sandbox *api.Sandbox, dir string, content string) []string {
	specs, err := sandbox.Deps.SerializableDeps.ParseYaml(content)
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
		violations = append(violations, checkRawParameter(sandbox, dir, entry)...)
	}

	return violations
}

// checkRawParameter reports the two contradictions one unparsed parameter
// entry may carry.
func checkRawParameter(sandbox *api.Sandbox, dir string, entry *serializabledeps.SerializableObject) []string {
	if !rawBool(entry, "required") {
		return nil
	}

	label := rawString(entry, "id")

	var violations []string
	if item, _ := entry.GetObjectItem("default"); item != nil && !item.IsNull() {
		violations = append(violations, routeViolation(dir,
			"declares the parameter "+label+" as both required and defaulted; the default already covers its absence"))
	}
	if rawString(entry, "type") == "boolean" || rawString(entry, "type") == "bool" {
		violations = append(violations, routeViolation(dir,
			"declares the boolean parameter "+label+" as required; its absence already means false"))
	}
	return violations
}

func rawString(entry *serializabledeps.SerializableObject, key string) string {
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

func rawBool(entry *serializabledeps.SerializableObject, key string) bool {
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

// routeViolation words one violation the same way for every rule.
func routeViolation(dir string, reason string) string {
	return dir + " " + reason
}
