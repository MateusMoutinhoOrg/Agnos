package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
)

// The matcher below is the generated sandbox/internal/server/route/IsActionable.go
// read against a route.yaml instead of an api.Route, so explain-route can say
// which routes a request reaches without a server. The two are kept one for
// one — a change to one is a change to the other, and the explain-route
// example holds them together.

// RouteRequest is one request as explain-route is handed it: the method, the
// path without its query string, and every value a parameter may be read from.
type RouteRequest struct {
	Method  string
	Path    string
	Query   map[string][]string
	Headers map[string]string
	Cookies map[string]string
}

// RouteMatch is what one route says about one request: whether it runs, and
// why not when it does not. MethodMismatch reports a route whose path matched
// under another method — what tells a 405 from a 404. Failure is the 400 its
// binding would answer with, "" when every value binds.
type RouteMatch struct {
	Runs           bool
	MethodMismatch bool
	Reason         string
	Failure        string
}

// MatchRouteRequest runs one route's declaration against one request.
func MatchRouteRequest(sandbox *api.Sandbox, conf *routeconf.RouteConf, request RouteRequest) RouteMatch {
	segments := []string{}
	for _, segment := range sandbox.Deps.Stringsdeps.Split(request.Path, "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}

	if conf.HasSegments && len(segments) != conf.Segments {
		return RouteMatch{Reason: sandbox.Deps.Std.Sprintf("the route takes %d segments, the path has %d", conf.Segments, len(segments))}
	}

	for _, path := range conf.Paths {
		text, ok := routePathSlice(sandbox, segments, path)
		if !ok {
			return RouteMatch{Reason: sandbox.Deps.Std.Sprintf("path %s: the request has no segment %d", path.Id, path.Start)}
		}
		if !routePathConverts(sandbox, path, text) {
			return RouteMatch{Reason: sandbox.Deps.Std.Sprintf("path %s: %q is not a valid %s", path.Id, sandbox.Deps.Stringsdeps.TrimPrefix(text, "/"), path.Type)}
		}
		if path.Trigger.Exists && !MatchRouteTrigger(sandbox, path.Trigger, text, true) {
			return RouteMatch{Reason: sandbox.Deps.Std.Sprintf("path %s: %q fails %s", path.Id, text, describeTrigger(sandbox, path.Trigger))}
		}
	}

	if !routeAccepts(conf, request.Method) {
		return RouteMatch{MethodMismatch: true, Reason: sandbox.Deps.Std.Sprintf("the method %s is not one of %s", request.Method, sandbox.Deps.Stringsdeps.Join(conf.Methods, ", "))}
	}

	for _, parameter := range conf.Parameters {
		if !parameter.Trigger.Exists {
			continue
		}
		values := routeParameterValues(sandbox, parameter, request)
		if len(values) == 0 {
			return RouteMatch{Reason: sandbox.Deps.Std.Sprintf("parameter %s: absent, and it declares a trigger", parameter.Id)}
		}
		if !MatchRouteTrigger(sandbox, parameter.Trigger, values[0], false) {
			return RouteMatch{Reason: sandbox.Deps.Std.Sprintf("parameter %s: %q fails %s", parameter.Id, values[0], describeTrigger(sandbox, parameter.Trigger))}
		}
	}

	match := RouteMatch{Runs: true}
	for _, parameter := range conf.Parameters {
		values := routeParameterValues(sandbox, parameter, request)
		if len(values) == 0 {
			if parameter.Required {
				match.Failure = sandbox.Deps.Std.Sprintf("400, required parameter %q is missing", parameter.Key)
				return match
			}
			continue
		}
		if !routeParameterConverts(sandbox, parameter, values) {
			match.Failure = sandbox.Deps.Std.Sprintf("400, parameter %q is not a valid %s", parameter.Key, parameter.Type)
			return match
		}
	}
	return match
}

// MatchRouteTrigger is MatchTrigger of the generated IsActionable.go.
func MatchRouteTrigger(sandbox *api.Sandbox, trigger routeconf.Trigger, text string, segmented bool) bool {
	value := trigger.Value
	if trigger.IgnoreCase && trigger.Type != "regex" {
		value = sandbox.Deps.Stringsdeps.ToLower(value)
		text = sandbox.Deps.Stringsdeps.ToLower(text)
	}

	matched := text == value
	switch trigger.Type {
	case "prefix":
		if !segmented {
			matched = sandbox.Deps.Stringsdeps.HasPrefix(text, value)
			break
		}
		value = sandbox.Deps.Stringsdeps.TrimSuffix(value, "/")
		matched = value == "" || text == value || sandbox.Deps.Stringsdeps.HasPrefix(text, value+"/")
	case "text-prefix":
		matched = sandbox.Deps.Stringsdeps.HasPrefix(text, value)
	case "suffix":
		matched = sandbox.Deps.Stringsdeps.HasSuffix(text, value)
	case "regex":
		if trigger.IgnoreCase {
			value = "(?i)" + value
		}
		ok, err := sandbox.Deps.Stringsdeps.MatchPattern(value, text)
		matched = err == nil && ok
	}
	return matched != trigger.Negate
}

// describeTrigger words one trigger the way explain-route prints it.
func describeTrigger(sandbox *api.Sandbox, trigger routeconf.Trigger) string {
	text := sandbox.Deps.Std.Sprintf("%s %q", trigger.Type, trigger.Value)
	if trigger.IgnoreCase {
		text += " ignoring case"
	}
	if trigger.Negate {
		text = "not " + text
	}
	return text
}

// routePathSlice is PathSlice of the generated IsActionable.go.
func routePathSlice(sandbox *api.Sandbox, segments []string, path routeconf.Path) (string, bool) {
	end := path.End
	if end < 0 {
		end = len(segments) - 1
	}
	if len(segments) == 0 && path.Start == 0 && path.End < 0 {
		return "/", true
	}
	if path.Start < 0 || path.Start > end || end >= len(segments) {
		return "", false
	}
	return "/" + sandbox.Deps.Stringsdeps.Join(segments[path.Start:end+1], "/"), true
}

// routeUuidPattern is uuidPattern of the generated IsActionable.go.
const routeUuidPattern = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`

// routePathConverts is PathValue of the generated IsActionable.go, reduced to
// whether the slice converts.
func routePathConverts(sandbox *api.Sandbox, path routeconf.Path, text string) bool {
	segment := sandbox.Deps.Stringsdeps.TrimPrefix(text, "/")
	switch path.Type {
	case "integer":
		_, err := sandbox.Deps.Stringsdeps.Atoi(segment)
		return err == nil
	case "number":
		_, err := sandbox.Deps.Stringsdeps.ParseFloat(segment, 64)
		return err == nil
	case "uuid":
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(routeUuidPattern, segment)
		return err == nil && matched
	}
	return true
}

// routeAccepts is acceptsMethod of the generated IsActionable.go.
func routeAccepts(conf *routeconf.RouteConf, method string) bool {
	for _, accepted := range conf.Methods {
		if accepted == method || accepted == routeconf.AnyMethod {
			return true
		}
	}
	return false
}

// routeParameterValues is ParameterValues of the generated IsActionable.go.
func routeParameterValues(sandbox *api.Sandbox, parameter routeconf.Parameter, request RouteRequest) []string {
	array := parameter.Type == "string-array" || parameter.Type == "integer-array"
	for _, font := range parameter.Fonts {
		raws := []string{}
		switch font {
		case "header":
			raw := request.Headers[sandbox.Deps.Stringsdeps.ToLower(parameter.Key)]
			if array {
				raws = sandbox.Deps.Stringsdeps.Split(raw, ",")
			} else {
				raws = []string{raw}
			}
		case "query":
			raws = request.Query[parameter.Key]
			if !array && len(raws) > 1 {
				raws = raws[:1]
			}
		case "cookie":
			raws = []string{request.Cookies[parameter.Key]}
		}

		values := []string{}
		for _, raw := range raws {
			if value := sandbox.Deps.Stringsdeps.TrimSpace(raw); value != "" {
				values = append(values, value)
			}
		}
		if len(values) > 0 {
			return values
		}
	}
	return []string{}
}

// routeParameterConverts reports whether the values a parameter brought bind
// to its type, as parseValue of the generated RequestHandler.go reads them.
func routeParameterConverts(sandbox *api.Sandbox, parameter routeconf.Parameter, values []string) bool {
	for _, value := range values {
		switch parameter.Type {
		case "integer", "integer-array":
			if _, err := sandbox.Deps.Stringsdeps.Atoi(value); err != nil {
				return false
			}
		case "number":
			if _, err := sandbox.Deps.Stringsdeps.ParseFloat(value, 64); err != nil {
				return false
			}
		case "boolean":
			lower := sandbox.Deps.Stringsdeps.ToLower(value)
			if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
				return false
			}
		}
		if parameter.Type != "string-array" && parameter.Type != "integer-array" {
			return true
		}
	}
	return true
}
