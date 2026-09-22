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
// its own NewRoute and collected by sandbox/internal/server/new.go, so this file is
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
// runs every route of Server.Routes the request matches, in the order the
// collector put them — lowest `priority` first, and within one rung the route
// fixing the most literal segments first, so a route on "/" can never swallow
// one on "/home".
//
// More than one route may match, which is what a chain is: each one runs in
// turn until one of them **sets a status**. A handler that writes no status has
// not answered, so the next route runs — that is the whole of what makes a
// middleware a middleware. A handler that returns an error without answering
// ends the chain too, through HandleServerError.
//
// Nothing here writes a response itself. Every way a request can end without a
// route answering it — nothing matched, matched under another method, a field
// that will not bind, a panic — is raised through routeio.Fail and answered by
// one of the project's own Handle* files.
func dispatch(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response) {
	tracked, answered := routeio.Tracked(response)
	defer recoverRoute(sandbox, request, tracked, answered)

	segments := splitPath(sandbox, request.GetPath())
	method := request.GetMethod()
	method_mismatch := false
	ran := false

	for index := range sandbox.Server.Routes {
		declared := &sandbox.Server.Routes[index]
		if !matchPath(sandbox, declared, segments) {
			continue
		}
		if declared.Method != method {
			method_mismatch = true
			continue
		}
		if !matchFields(sandbox, declared, request) {
			continue
		}

		ran = true
		if runRoute(sandbox, declared, request, tracked, segments, answered) {
			return
		}
	}

	// A path some route answers under another method is the one failure the
	// chain can tell apart from a path nothing knows at all — and only when
	// no route ran, since a chain that ran and wrote nothing is the project
	// declining to answer, not the method being wrong.
	if !ran && method_mismatch {
		failRequest(sandbox, request, tracked, api.StatusMethodNotAllowed)
		return
	}

	failRequest(sandbox, request, tracked, api.StatusNotFound)
}

// matchPath reports whether a request path binds to one route's declaration:
// every literal segment has to be the one declared, a capture takes whatever
// sits in its place, and the segment count has to be exact — unless the route
// ends in something that takes the rest of the path. There are two of those: an
// array capture, which binds every segment left, and a prefix trigger, which
// matches the segments it spells and leaves everything after them unread.
func matchPath(sandbox *api.Sandbox, route *api.Route, segments []string) bool {
	position := 0

	for _, segment := range route.Paths {
		if segment.Field == nil {
			literals := literalSegments(sandbox, segment)

			if segment.StartsWith {
				// "/" spells no segment of its own, so it is the prefix
				// every path begins with.
				for _, literal := range literals {
					if position >= len(segments) || segments[position] != literal {
						return false
					}
					position++
				}
				return true
			}

			// The root identifier "/" fixes the empty path: it names no
			// segment of its own, so there is nothing to match on.
			for _, literal := range literals {
				if position >= len(segments) || segments[position] != literal {
					return false
				}
				position++
			}
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

// matchFields reports whether the request satisfies every value condition the
// route's headers and query parameters declare. A field carrying an
// `identifier` or a `starts-with-identifier` is part of what the route matches
// on: the request has to bring that value for the route to run at all. A field
// carrying neither is bound later and never matched on here.
func matchFields(sandbox *api.Sandbox, route *api.Route, request serverdeps.Request) bool {
	for _, field := range route.Headers {
		if !matchValue(sandbox, field, request.GetHeader(field.Id)) {
			return false
		}
	}
	for _, field := range route.Params {
		if !matchValue(sandbox, field, request.GetQueryParam(field.Id)) {
			return false
		}
	}
	return true
}

// matchValue tests one raw request value against the conditions its field
// declares. The two are never declared together, and a field declaring neither
// matches whatever the request brought, absent included.
func matchValue(sandbox *api.Sandbox, field api.RouteField, raw string) bool {
	if field.Identifier != "" && raw != field.Identifier {
		return false
	}
	if field.StartsWith != "" && !sandbox.Deps.Stringsdeps.HasPrefix(raw, field.StartsWith) {
		return false
	}
	return true
}

// literalSegments slices one trigger's identifier into the path segments it
// fixes. An exact identifier spells one, "/" spells none, and a prefix trigger
// spells as many as it names.
func literalSegments(sandbox *api.Sandbox, segment api.RoutePath) []string {
	var literals []string
	for _, literal := range sandbox.Deps.Stringsdeps.Split(segment.Identifier, "/") {
		if literal != "" {
			literals = append(literals, literal)
		}
	}
	return literals
}

// runRoute binds one request to one route's declaration and hands the result to
// its handler. The order is the order route.yaml reads in: the captured
// segments, then the headers, then the query parameters, then what none of them
// brought, then the body's own declaration.
//
// What it binds is a copy of the declaration, never the declaration on
// Server.Routes: two requests in flight hold their own Items and their own
// request.
//
// It reports whether the chain is over — because this route answered, or
// because a failure was raised that ends the request whatever else would have
// matched.
func runRoute(sandbox *api.Sandbox, declared *api.Route, request serverdeps.Request, response serverdeps.Response, segments []string, answered func() bool) bool {
	route := api.BindRoute(declared)
	route.Request = request
	route.Response = response

	if !bindPaths(sandbox, route, segments) {
		return true
	}
	if !bindFields(sandbox, route, route.Headers, subjectHeader, func(field api.RouteField) []string {
		return headerValues(request, field)
	}) {
		return true
	}
	if !bindFields(sandbox, route, route.Params, subjectParam, func(field api.RouteField) []string {
		return paramValues(request, field)
	}) {
		return true
	}
	if !bindMissing(sandbox, route) {
		return true
	}
	if !checkBody(sandbox, route, request) {
		return true
	}

	err := route.Handler(route)

	// The status is what ends the chain, so a handler that answered *and*
	// returned something has answered: what it returned is reported and goes
	// no further.
	if answered() {
		if err != nil {
			sandbox.Deps.Std.Log("route %s: %s \n", route.Name, err.Error())
		}
		return true
	}

	if err != nil {
		routeio.FailWithCause(sandbox, route, api.StatusFailure, "", "", err.Error())
		return true
	}

	return false
}

// bindPaths binds the captured segments of a matched path. A capture is always
// filled — the route only matched because that segment was there — so nothing
// here falls back to a default.
func bindPaths(sandbox *api.Sandbox, route *api.Route, segments []string) bool {
	position := 0

	for _, segment := range route.Paths {
		field := segment.Field
		if field == nil {
			position += len(literalSegments(sandbox, segment))
			continue
		}

		if field.Array {
			for _, raw := range segments[position:] {
				value, ok := parseValue(sandbox, route, subjectPath, *field, raw)
				if !ok {
					return false
				}
				route.Items[field.Id] = append(route.Items[field.Id], value)
			}
			return true
		}

		value, ok := parseValue(sandbox, route, subjectPath, *field, segments[position])
		if !ok {
			return false
		}
		if !inRange(sandbox, route, subjectPath, *field, value) {
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
func bindFields(sandbox *api.Sandbox, route *api.Route, fields []api.RouteField, subject string, values func(field api.RouteField) []string) bool {
	for _, field := range fields {
		if len(route.Items[field.Id]) > 0 {
			continue
		}

		raws := values(field)

		if field.Array {
			for _, raw := range raws {
				value, ok := parseValue(sandbox, route, subject, field, raw)
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

		value, ok := parseValue(sandbox, route, subject, field, raws[0])
		if !ok {
			return false
		}
		if !inRange(sandbox, route, subject, field, value) {
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
func bindMissing(sandbox *api.Sandbox, route *api.Route) bool {
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
				routeio.Fail(sandbox, route, api.StatusBadRequest, field.Id,
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
func checkBody(sandbox *api.Sandbox, route *api.Route, request serverdeps.Request) bool {
	if route.Body.Type == "" || route.Body.Type == bodyNone {
		return true
	}

	if route.Body.ContentType != "" {
		content_type := request.GetHeader("Content-Type")
		if content_type != "" && !sandbox.Deps.Stringsdeps.HasPrefix(content_type, route.Body.ContentType) {
			routeio.Fail(sandbox, route, api.StatusUnsupportedMedia, "",
				sandbox.Deps.Std.Sprintf("this route accepts a %s body", route.Body.ContentType))
			return false
		}
	}

	if raw := request.GetHeader("Content-Length"); raw != "" {
		declared, err := sandbox.Deps.Stringsdeps.Atoi(raw)
		if err == nil && declared > route.Body.MaxBytes {
			routeio.Fail(sandbox, route, api.StatusPayloadTooLarge, "",
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

// failRequest raises a failure that belongs to no route — nothing matched the
// path, or nothing matched it under this method. The handler still gets an
// api.Route, carrying the request and the response and nothing else, so every
// Handle* file reads the same whichever failure brought it there.
//
// It carries no message on purpose. There is nothing to say about these two
// beyond the status, so the wording is the project's: routeio.FailureOf fills
// in the one the Handle* file spells, which is what makes editing that file
// change what the server says. The failures that do carry information — a
// field that would not bind, a body the schema rejected — pass their message
// through instead.
func failRequest(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response, status int) {
	route := api.NewRoute()
	route.Request = request
	route.Response = response

	routeio.Fail(sandbox, route, status, "", "")
}

// recoverRoute turns a panicking handler into one answered request, so a single
// bad route cannot take the process down with it. A handler that panicked after
// writing its status has already answered: the panic is reported and nothing is
// written over it.
func recoverRoute(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response, answered func() bool) {
	failure := recover()
	if failure == nil {
		return
	}

	sandbox.Deps.Std.Error("route panicked: %v\n", failure)
	if answered() {
		return
	}

	route := api.NewRoute()
	route.Request = request
	route.Response = response

	routeio.FailWithCause(sandbox, route, api.StatusFailure, "", "",
		sandbox.Deps.Std.Sprintf("%v", failure))
}

// parseValue converts one raw request value to the type its declaration names,
// answering 400 in the server's own words rather than the conversion library's.
// Every raw value is already a string, so only the other three convert.
func parseValue(sandbox *api.Sandbox, route *api.Route, subject string, field api.RouteField, raw string) (any, bool) {
	switch field.Type {
	case typeInt:
		value, err := sandbox.Deps.Stringsdeps.Atoi(raw)
		if err != nil {
			routeio.Fail(sandbox, route, api.StatusBadRequest, field.Id,
				sandbox.Deps.Std.Sprintf("%s '%s' is not a valid integer", subject, field.Id))
			return nil, false
		}
		return value, true
	case typeFloat:
		value, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
		if err != nil {
			routeio.Fail(sandbox, route, api.StatusBadRequest, field.Id,
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
		routeio.Fail(sandbox, route, api.StatusBadRequest, field.Id,
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
func inRange(sandbox *api.Sandbox, route *api.Route, subject string, field api.RouteField, value any) bool {
	if !field.HasMin && !field.HasMax {
		return true
	}

	number, ok := numberOf(value)
	if !ok {
		return true
	}

	if field.HasMin && number < field.Min {
		routeio.Fail(sandbox, route, api.StatusBadRequest, field.Id,
			sandbox.Deps.Std.Sprintf("%s '%s' must be >= %s", subject, field.Id, numberLabel(sandbox, field.Type, field.Min)))
		return false
	}
	if field.HasMax && number > field.Max {
		routeio.Fail(sandbox, route, api.StatusBadRequest, field.Id,
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
