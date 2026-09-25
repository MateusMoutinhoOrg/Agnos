package route

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeio"
)

// entriesArgument is the position of the Entries pointer among the parameters
// of an InternalPurehandler: func(route, entries, response) error.
const entriesArgument = 1

// entriesTag is the struct tag an Entries field names what it is bound to by.
const entriesTag = "id"

// fullRouteId is the id every Entries carries the whole request path under.
const fullRouteId = "FullRoute"

// datetimePattern is what a `datetime` parameter has to read as: RFC 3339.
const datetimePattern = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})$`

// RequestHandler binds the request of one bound route onto a fresh Entries and
// runs the route's InternalPurehandler with it. The Entries type is the route
// package's own, so it is built, filled and called through Deps.Reflectdeps:
// every field is filled by its `id` tag — FullRoute, one per entry of Paths,
// one per parameter.
//
// The order is the order route.yaml reads in: the path slices, then the
// parameters, then the body's own declaration, then the response type. A
// value that will not bind is raised through routeio.Fail and the handler never
// runs. What it returns is what the handler returned.
func RequestHandler(sandbox *api.Sandbox, route *api.Route) error {
	request := routeio.RequestOf(route)
	response := routeio.ResponseOf(route)

	entries := sandbox.Deps.Reflectdeps.NewIn(route.InternalPurehandler, entriesArgument)
	if entries == nil || sandbox.Deps.Reflectdeps.NumField(entries) < 0 {
		return sandbox.Deps.Std.Errorf("route %s: InternalPurehandler is not a func(route *api.Route, entries *Entries, response *serverdeps.Response) error", route.Name)
	}

	values, ok, err := bindValues(sandbox, route, request)
	if !ok {
		return err
	}
	if ok, err := checkBody(sandbox, route, request); !ok {
		return err
	}

	for index := 0; index < sandbox.Deps.Reflectdeps.NumField(entries); index++ {
		id := sandbox.Deps.Reflectdeps.FieldTag(entries, index, entriesTag)
		value, has := values[id]
		if id == "" || !has {
			continue
		}
		if err := sandbox.Deps.Reflectdeps.SetField(entries, index, value); err != nil {
			return sandbox.Deps.Std.Errorf("route %s: Entries.%s: %s", route.Name,
				sandbox.Deps.Reflectdeps.FieldName(entries, index), err.Error())
		}
	}

	if route.ResponseType != "" {
		response.SetHeader("Content-Type", route.ResponseType)
	}

	out := sandbox.Deps.Reflectdeps.Call(route.InternalPurehandler, []any{route, entries, &response})
	if len(out) == 1 && out[0] != nil {
		if handler_error, is := out[0].(error); is {
			return handler_error
		}
	}
	return nil
}

// bindValues reads every value the route declares off the request, keyed by
// the id its Entries field is tagged with. A required parameter the request
// does not bring, or a value that will not convert, is raised through
// routeio.Fail: it reports false, with what that failure returned.
func bindValues(sandbox *api.Sandbox, route *api.Route, request serverdeps.Request) (map[string]any, bool, error) {
	values := map[string]any{fullRouteId: request.GetPath()}

	segments := SplitPath(sandbox, request.GetPath())
	for _, path := range route.Paths {
		text, _ := PathSlice(sandbox, segments, path)
		values[path.Id] = text
	}

	for _, parameter := range route.Parameters {
		raws := ParameterValues(sandbox, request, parameter)

		if len(raws) == 0 {
			if parameter.Required {
				return nil, false, routeio.Fail(sandbox, route, api.StatusBadRequest, parameter.Key,
					sandbox.Deps.Std.Sprintf("required parameter '%s' is missing", parameter.Key))
			}
			if !parameter.HasDefault {
				continue
			}
			raws = []string{parameter.Default}
		}

		value, ok := parseValue(sandbox, parameter, raws)
		if !ok {
			return nil, false, routeio.Fail(sandbox, route, api.StatusBadRequest, parameter.Key,
				sandbox.Deps.Std.Sprintf("parameter '%s' %s", parameter.Key, typeMessage(parameter.Type)))
		}
		values[parameter.Id] = value
	}

	return values, true, nil
}

// parseValue converts the raw values of one parameter to the type it
// declares, reporting false when they will not.
func parseValue(sandbox *api.Sandbox, parameter api.Parameter, raws []string) (any, bool) {
	raw := raws[0]

	switch parameter.Type {
	case api.NumberType:
		value, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
		return value, err == nil
	case api.BooleanType:
		switch sandbox.Deps.Stringsdeps.ToLower(raw) {
		case "true", "1":
			return true, true
		case "false", "0":
			return false, true
		}
		return false, false
	case api.DateTimeType:
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(datetimePattern, raw)
		return raw, err == nil && matched
	case api.StringArrayType:
		return raws, true
	}
	return raw, true
}

// typeMessage is how a value that will not convert is reported, in the
// server's own words rather than the conversion library's.
func typeMessage(kind api.ParameterType) string {
	switch kind {
	case api.NumberType:
		return "is not a valid number"
	case api.BooleanType:
		return "must be true or false"
	case api.DateTimeType:
		return "is not an RFC 3339 date-time"
	}
	return "is not valid"
}

// checkBody enforces what the body's declaration settles before a byte of it is
// read: the media type and the length the request itself declares. The body is
// the handler's to ask for, through the ReadBody its own entries.go generates.
func checkBody(sandbox *api.Sandbox, route *api.Route, request serverdeps.Request) (bool, error) {
	if route.Body.Type == "" || route.Body.Type == "none" {
		return true, nil
	}

	if route.Body.ContentType != "" {
		content_type := request.GetHeader("Content-Type")
		if content_type != "" && !sandbox.Deps.Stringsdeps.HasPrefix(content_type, route.Body.ContentType) {
			return false, routeio.Fail(sandbox, route, api.StatusUnsupportedMedia, "",
				sandbox.Deps.Std.Sprintf("this route accepts a %s body", route.Body.ContentType))
		}
	}

	if raw := request.GetHeader("Content-Length"); raw != "" {
		declared, err := sandbox.Deps.Stringsdeps.Atoi(raw)
		if err == nil && declared > route.Body.MaxBytes {
			return false, routeio.Fail(sandbox, route, api.StatusPayloadTooLarge, "",
				sandbox.Deps.Std.Sprintf("the request body is larger than %s bytes",
					sandbox.Deps.Stringsdeps.FormatInt(int64(route.Body.MaxBytes), 10)))
		}
	}

	return true, nil
}
