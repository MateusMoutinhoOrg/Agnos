package server

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeio"
)

// The declared types a route field may carry, as route.yaml spells them.
const (
	typeBoolean = "boolean"
	typeInt     = "int"
	typeFloat   = "float"
)

// The body types a route may declare. A route declaring none reads nothing and
// is held to no content type or length.
const bodyNone = "none"

// subject* is how a binding error names where a value came from.
const (
	subjectPath   = "path segment"
	subjectHeader = "header"
	subjectParam  = "query parameter"
)

// ServerMain opens the port through sandbox.Deps.Serverdeps — which routes
// nothing — and hands every request to dispatch, whatever its method or path.
// Nothing here is generated per route: every route is one declaration built by
// its own NewRoute and collected by sandbox/binds/server.go, so this file is
// the same in every project.
func ServerMain(sandbox *api.Sandbox, props api.ServeProps) error {
	server := sandbox.Deps.Serverdeps.NewServer(serverdeps.ServerProps{
		Addr:           props.Addr,
		ReadTimeoutMs:  props.ReadTimeoutMs,
		WriteTimeoutMs: props.WriteTimeoutMs,
		Handler: func(request serverdeps.Request, response serverdeps.Response) {
			dispatch(sandbox, request, response)
		},
	})

	sandbox.Deps.Std.Log("server listening on %s \n", props.Addr)
	return server.Listen()
}

// dispatch is the whole routing layer: it slices the path into segments and
// tests every route of sandbox.Routes against them, in the order the collector
// put them — the route fixing the most literal segments first, so a route on
// "/" can never swallow one on "/home".
//
// Everything but the body is settled here: a path nothing matches is 404, a
// path matched under another method is 405, and a field that will not bind is
// 400. A RouteHandler never decides any of those.
func dispatch(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response) {
	defer recoverRoute(sandbox, response)

	segments := splitPath(sandbox, request.GetPath())
	method := request.GetMethod()
	path_matched := false

	for _, declared := range sandbox.Routes {
		if !matchRoute(sandbox, declared, segments) {
			continue
		}
		path_matched = true
		if declared.Method != method {
			continue
		}
		runRoute(sandbox, declared, request, response, segments)
		return
	}

	if path_matched {
		routeio.WriteError(sandbox, response, api.StatusMethodNotAllowed, "", "method not allowed")
		return
	}

	routeio.WriteError(sandbox, response, api.StatusNotFound, "", "route not found")
}

// matchRoute reports whether a request path binds to one route's declaration:
// every literal segment has to be the one declared, a capture takes whatever
// sits in its place, and the segment count has to be exact — unless the last
// capture is an array, which takes every segment left and so wants strictly
// more than the route fixes.
func matchRoute(sandbox *api.Sandbox, route *api.Route, segments []string) bool {
	position := 0

	for _, segment := range route.Paths {
		if segment.Field == nil {
			// The root identifier "/" fixes the empty path: it names no
			// segment of its own, so there is nothing to match on.
			literal := sandbox.Deps.Stringsdeps.Trim(segment.Identifier, "/")
			if literal == "" {
				continue
			}
			if position >= len(segments) || segments[position] != literal {
				return false
			}
			position++
			continue
		}

		if segment.Field.Array {
			// The capture taking the rest of the path is required like any
			// other, so it wants one segment at least.
			return len(segments) > position
		}

		if position >= len(segments) {
			return false
		}
		position++
	}

	return len(segments) == position
}

// runRoute binds one request to one route's declaration and hands the result to
// its handler. The order is the order route.yaml reads in: the captured
// segments, then the headers, then the query parameters, then what none of them
// brought, then the body's own declaration.
//
// The route it binds is a fresh instance, minted by the declaration's own New:
// two requests in flight hold their own Items and their own request.
func runRoute(sandbox *api.Sandbox, declared *api.Route, request serverdeps.Request, response serverdeps.Response, segments []string) {
	route := declared.New()
	route.Request = request
	route.Response = response

	if !bindPaths(sandbox, route, response, segments) {
		return
	}
	if !bindFields(sandbox, route, response, route.Headers, subjectHeader, func(field api.RouteField) []string {
		return headerValues(request, field)
	}) {
		return
	}
	if !bindFields(sandbox, route, response, route.Params, subjectParam, func(field api.RouteField) []string {
		return paramValues(request, field)
	}) {
		return
	}
	if !bindMissing(sandbox, route, response) {
		return
	}
	if !checkBody(sandbox, route, request, response) {
		return
	}

	if status := route.Handler(); status == 0 {
		routeio.WriteError(sandbox, response, api.StatusFailure, "", "the route handler returned no status")
	}
}

// bindPaths binds the captured segments of a matched path. A capture is always
// filled — the route only matched because that segment was there — so nothing
// here falls back to a default.
func bindPaths(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response, segments []string) bool {
	position := 0

	for _, segment := range route.Paths {
		field := segment.Field
		if field == nil {
			if sandbox.Deps.Stringsdeps.Trim(segment.Identifier, "/") != "" {
				position++
			}
			continue
		}

		if field.Array {
			for _, raw := range segments[position:] {
				value, ok := parseValue(sandbox, response, subjectPath, *field, raw)
				if !ok {
					return false
				}
				route.Items[field.Id] = append(route.Items[field.Id], value)
			}
			return true
		}

		value, ok := parseValue(sandbox, response, subjectPath, *field, segments[position])
		if !ok {
			return false
		}
		if !inRange(sandbox, response, subjectPath, *field, value) {
			return false
		}
		route.Items[field.Id] = []any{value}
		position++
	}

	return true
}

// bindFields binds one origin's declared fields into route.Items, reading each
// raw value through values. A name may be declared in more than one origin: the
// field is bound once and filled by the first origin, in paths -> headers ->
// params order, that brings a value, so one already bound is left alone.
func bindFields(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response, fields []api.RouteField, subject string, values func(field api.RouteField) []string) bool {
	for _, field := range fields {
		if len(route.Items[field.Id]) > 0 {
			continue
		}

		raws := values(field)

		if field.Array {
			for _, raw := range raws {
				value, ok := parseValue(sandbox, response, subject, field, raw)
				if !ok {
					return false
				}
				route.Items[field.Id] = append(route.Items[field.Id], value)
			}
			continue
		}

		if len(raws) == 0 || raws[0] == "" {
			continue
		}

		value, ok := parseValue(sandbox, response, subject, field, raws[0])
		if !ok {
			return false
		}
		if !inRange(sandbox, response, subject, field, value) {
			return false
		}
		route.Items[field.Id] = []any{value}
	}

	return true
}

// bindMissing settles every declared header and query parameter the request
// brought nothing for: a required one is 400 whichever origin declared it, and
// the rest fall back to their declared default. The required pass runs first,
// so a name declared required in one origin and defaulted in another is still
// reported rather than quietly defaulted.
func bindMissing(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) bool {
	origins := []struct {
		fields  []api.RouteField
		subject string
	}{
		{route.Headers, subjectHeader},
		{route.Params, subjectParam},
	}

	for _, origin := range origins {
		for _, field := range origin.fields {
			if field.Required && len(route.Items[field.Id]) == 0 {
				routeio.WriteError(sandbox, response, api.StatusBadRequest, field.Id,
					sandbox.Deps.Std.Sprintf("required %s '%s' is missing", origin.subject, field.Id))
				return false
			}
		}
	}

	for _, origin := range origins {
		for _, field := range origin.fields {
			if field.HasDefault && len(route.Items[field.Id]) == 0 {
				route.Items[field.Id] = []any{defaultValue(sandbox, field.Type, field.Default)}
			}
		}
	}

	return true
}

// checkBody enforces what the body's declaration settles before a byte of it is
// read: the media type and the length the request itself declares. The body is
// the handler's to ask for, through the ReadBody its own package generates.
func checkBody(sandbox *api.Sandbox, route *api.Route, request serverdeps.Request, response serverdeps.Response) bool {
	if route.Body.Type == "" || route.Body.Type == bodyNone {
		return true
	}

	if route.Body.ContentType != "" {
		content_type := request.GetHeader("Content-Type")
		if content_type != "" && !sandbox.Deps.Stringsdeps.HasPrefix(content_type, route.Body.ContentType) {
			routeio.WriteError(sandbox, response, api.StatusUnsupportedMedia, "",
				sandbox.Deps.Std.Sprintf("this route accepts a %s body", route.Body.ContentType))
			return false
		}
	}

	if raw := request.GetHeader("Content-Length"); raw != "" {
		declared, err := sandbox.Deps.Stringsdeps.Atoi(raw)
		if err == nil && declared > route.Body.MaxBytes {
			routeio.WriteError(sandbox, response, api.StatusPayloadTooLarge, "",
				sandbox.Deps.Std.Sprintf("the request body is larger than %s bytes",
					sandbox.Deps.Stringsdeps.FormatInt(int64(route.Body.MaxBytes), 10)))
			return false
		}
	}

	return true
}

// ─── request helpers ────────────────────────────────────────────────────────
//
// Every value bound above goes through these, so one spelling of every binding
// error is emitted for the whole server.

// headerValues is the raw values one declared header brings. A header is a
// single value whatever its declaration says: `array: true` is refused there.
func headerValues(request serverdeps.Request, field api.RouteField) []string {
	return []string{request.GetHeader(field.Id)}
}

// paramValues is the raw values one declared query parameter brings: every
// occurrence of the key for an array, the first one otherwise.
func paramValues(request serverdeps.Request, field api.RouteField) []string {
	if field.Array {
		return request.GetQueryAll(field.Id)
	}
	return []string{request.GetQueryParam(field.Id)}
}

// splitPath slices a raw request path into the segments a route is matched
// against, dropping the empty ones the leading and trailing slashes leave
// behind. The root path yields no segment at all.
func splitPath(sandbox *api.Sandbox, path string) []string {
	var segments []string
	for _, segment := range sandbox.Deps.Stringsdeps.Split(path, "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}
	return segments
}

// recoverRoute turns a panicking handler into one 500, so a single bad route
// cannot take the process down with it.
func recoverRoute(sandbox *api.Sandbox, response serverdeps.Response) {
	failure := recover()
	if failure == nil {
		return
	}
	sandbox.Deps.Std.Error("route panicked: %v\n", failure)
	routeio.WriteError(sandbox, response, api.StatusFailure, "", "internal error")
}

// parseValue converts one raw request value to the type its declaration names,
// answering 400 in the server's own words rather than the conversion library's.
// Every raw value is already a string, so only the other three convert.
func parseValue(sandbox *api.Sandbox, response serverdeps.Response, subject string, field api.RouteField, raw string) (any, bool) {
	switch field.Type {
	case typeInt:
		value, err := sandbox.Deps.Stringsdeps.Atoi(raw)
		if err != nil {
			routeio.WriteError(sandbox, response, api.StatusBadRequest, field.Id,
				sandbox.Deps.Std.Sprintf("%s '%s' is not a valid integer", subject, field.Id))
			return nil, false
		}
		return value, true
	case typeFloat:
		value, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
		if err != nil {
			routeio.WriteError(sandbox, response, api.StatusBadRequest, field.Id,
				sandbox.Deps.Std.Sprintf("%s '%s' is not a valid number", subject, field.Id))
			return nil, false
		}
		return value, true
	case typeBoolean:
		switch sandbox.Deps.Stringsdeps.ToLower(raw) {
		case "true", "1":
			return true, true
		case "false", "0":
			return false, true
		}
		routeio.WriteError(sandbox, response, api.StatusBadRequest, field.Id,
			sandbox.Deps.Std.Sprintf("%s '%s' must be true or false", subject, field.Id))
		return nil, false
	}
	return raw, true
}

// defaultValue is the declared default, spelled in route.yaml as text, in the
// type the declaration names.
func defaultValue(sandbox *api.Sandbox, kind string, raw string) any {
	switch kind {
	case typeBoolean:
		return raw == "true"
	case typeInt:
		value, err := sandbox.Deps.Stringsdeps.Atoi(raw)
		if err != nil {
			return 0
		}
		return value
	case typeFloat:
		value, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
		if err != nil {
			return float64(0)
		}
		return value
	}
	return raw
}

// inRange enforces the min/max bounds a numeric field declares. A value of any
// other type carries no bounds and passes untested.
func inRange(sandbox *api.Sandbox, response serverdeps.Response, subject string, field api.RouteField, value any) bool {
	if !field.HasMin && !field.HasMax {
		return true
	}

	number, ok := numberOf(value)
	if !ok {
		return true
	}

	if field.HasMin && number < field.Min {
		routeio.WriteError(sandbox, response, api.StatusBadRequest, field.Id,
			sandbox.Deps.Std.Sprintf("%s '%s' must be >= %s", subject, field.Id, numberLabel(sandbox, field.Type, field.Min)))
		return false
	}
	if field.HasMax && number > field.Max {
		routeio.WriteError(sandbox, response, api.StatusBadRequest, field.Id,
			sandbox.Deps.Std.Sprintf("%s '%s' must be <= %s", subject, field.Id, numberLabel(sandbox, field.Type, field.Max)))
		return false
	}

	return true
}

// numberOf reads a bound value as the number a bound is compared against.
func numberOf(value any) (float64, bool) {
	if number, ok := value.(int); ok {
		return float64(number), true
	}
	if number, ok := value.(float64); ok {
		return number, true
	}
	return 0, false
}

// numberLabel spells a bound the way its declaration does, so the message
// quotes what route.yaml says.
func numberLabel(sandbox *api.Sandbox, kind string, value float64) string {
	if kind == typeInt {
		return sandbox.Deps.Stringsdeps.FormatInt(int64(value), 10)
	}
	return sandbox.Deps.Stringsdeps.FormatFloat(value, 'g', -1, 64)
}
