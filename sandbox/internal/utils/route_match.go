package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
)

// The matcher below is the generated sandbox/internal/generated/server/route/IsActionable.go
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
		if path.Trigger.Exists && !MatchTrigger(sandbox, path.Trigger, text, true) {
			return RouteMatch{Reason: sandbox.Deps.Std.Sprintf("path %s: %q fails %s", path.Id, text, DescribeTrigger(sandbox, path.Trigger))}
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
		if !MatchTrigger(sandbox, parameter.Trigger, values[0], false) {
			return RouteMatch{Reason: sandbox.Deps.Std.Sprintf("parameter %s: %q fails %s", parameter.Id, values[0], DescribeTrigger(sandbox, parameter.Trigger))}
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

// RouteMiddlewareReach crosses one route with another that runs on a lower
// rung — a middleware in front of it — the way MiddlewareReach crosses two
// commands: whether it runs in front of every request of the route (its
// methods take the route's, every path trigger holds on the route's literal
// segments), may run (a trigger reads a segment the route captures, or a
// regex), or never does; and, when one of its parameters declares a trigger,
// the condition that adds.
func RouteMiddlewareReach(sandbox *api.Sandbox, middleware *routeconf.RouteConf, route *routeconf.RouteConf) (Reach, string) {
	if middleware.Priority >= route.Priority {
		return NoReach, ""
	}

	reach := Runs
	if !routeMethodsMeet(middleware, route) {
		return NoReach, ""
	}
	if !contains(middleware.Methods, routeconf.AnyMethod) && contains(route.Methods, routeconf.AnyMethod) {
		reach = MayRun
	}
	if middleware.HasSegments && (!route.HasSegments || route.Segments != middleware.Segments) {
		if route.HasSegments {
			return NoReach, ""
		}
		reach = MayRun
	}

	words := routeLiteralSegments(sandbox, route)
	for _, path := range middleware.Paths {
		if !path.Trigger.Exists {
			continue
		}
		switch routeTriggerReach(sandbox, path, words) {
		case NoReach:
			return NoReach, ""
		case MayRun:
			reach = MayRun
		}
	}

	conditions := []string{}
	for _, parameter := range middleware.Parameters {
		if parameter.Trigger.Exists {
			conditions = append(conditions, parameter.Key+" "+DescribeTrigger(sandbox, parameter.Trigger))
		}
	}
	if len(conditions) == 0 {
		return reach, ""
	}
	return reach, "only when " + sandbox.Deps.Stringsdeps.Join(conditions, " and ")
}

// routeMethodsMeet reports whether two routes share a method they answer.
func routeMethodsMeet(middleware *routeconf.RouteConf, route *routeconf.RouteConf) bool {
	if contains(middleware.Methods, routeconf.AnyMethod) || contains(route.Methods, routeconf.AnyMethod) {
		return true
	}
	for _, method := range middleware.Methods {
		if contains(route.Methods, method) {
			return true
		}
	}
	return false
}

// routeLiteralSegments is, segment by segment, the one segment a route's
// request path is known to carry there: what an equal trigger of its paths
// spells.
func routeLiteralSegments(sandbox *api.Sandbox, route *routeconf.RouteConf) map[int]string {
	words := map[int]string{}
	for _, path := range route.Paths {
		if !path.Trigger.Exists || path.Trigger.Negate || path.Trigger.Type != "equal" {
			continue
		}
		offset := 0
		for _, segment := range sandbox.Deps.Stringsdeps.Split(path.Trigger.Value, "/") {
			if segment == "" {
				continue
			}
			if _, known := words[path.Start+offset]; !known {
				words[path.Start+offset] = segment
			}
			offset++
		}
	}
	return words
}

// routeTriggerReach is triggerReach for a path of a middleware route, over the
// route's literal segments.
func routeTriggerReach(sandbox *api.Sandbox, path routeconf.Path, words map[int]string) Reach {
	trigger := path.Trigger
	if trigger.Type == "regex" {
		return MayRun
	}

	known := []string{}
	complete := path.End != routeconf.LastSegment
	for index := path.Start; path.End == routeconf.LastSegment || index <= path.End; index++ {
		word, has := words[index]
		if !has {
			complete = false
			break
		}
		known = append(known, word)
	}
	if path.End == routeconf.LastSegment {
		complete = false
	}

	text := "/" + sandbox.Deps.Stringsdeps.Join(known, "/")
	if complete {
		return boolReach(MatchTrigger(sandbox, trigger, text, true))
	}

	if trigger.Type == "prefix" && !trigger.Negate {
		value := []string{}
		for _, segment := range sandbox.Deps.Stringsdeps.Split(trigger.Value, "/") {
			if segment != "" {
				value = append(value, segment)
			}
		}
		if len(value) <= len(known) {
			return boolReach(MatchTrigger(sandbox, trigger, text, true))
		}
		for index, word := range known {
			if value[index] != word {
				return NoReach
			}
		}
	}
	return MayRun
}
