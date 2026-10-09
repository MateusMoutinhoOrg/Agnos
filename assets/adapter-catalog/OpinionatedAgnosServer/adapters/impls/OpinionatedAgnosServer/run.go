package opinionatedagnosserver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	opinionatedagnosserver "{{.Module}}/sandbox/deps/OpinionatedAgnosServer"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// entriesArgument is the position of the Input pointer among the parameters
// of a Handle: func(props, entries, response) error.
const entriesArgument = 1

// entriesTag is the struct tag an Input field names what it is bound to by.
const entriesTag = "id"

// fullRouteId is the id every Input carries the whole request path under.
const fullRouteId = "FullRoute"

// bodyId is the id the Input of a route declaring a body carries it under.
const bodyId = "Body"

// datetimePattern is what a `datetime` parameter has to read as: RFC 3339.
const datetimePattern = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})$`

// handle binds the request of one bound route onto a fresh Input and runs
// the route's Handle with it. The Input type is the route
// package's own, so it is built, filled and called by reflection: every field
// is filled by its `id` tag — FullRoute, one per entry of Paths, one per
// parameter, and the body.
//
// The order is the order route.yaml reads in: the path slices, then the
// parameters, then the body, then the response type. A value that will not
// bind is raised and the handler never runs. The handler is handed the
// request's shared RouteProps first: what a route earlier in the chain set
// there is what it reads. A failure it — or the route's ReadBody — returns,
// built by Fail, is raised on this route; any other error is returned as it
// is.
func (server_run *run) handle(route *opinionatedagnosserver.Route) error {
	request := route.Request
	response := route.Response

	entries := newIn(route.Handle, entriesArgument)
	if entries == nil || numField(entries) < 0 {
		return fmt.Errorf("route %s: Handle is not a func(props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error", route.Name)
	}

	values, ok, err := server_run.bindValues(route, request)
	if !ok {
		return err
	}
	if ok, err := server_run.checkBody(route, request); !ok {
		return err
	}
	if route.ReadBody != nil {
		body, err := route.ReadBody(route)
		if err != nil {
			if failure, is := err.(*opinionatedagnosserver.RouteFailure); is {
				return server_run.raiseFailure(route, failure)
			}
			return err
		}
		values[bodyId] = body
	}

	for index := 0; index < numField(entries); index++ {
		id := fieldTag(entries, index, entriesTag)
		value, has := values[id]
		if id == "" || !has {
			continue
		}
		if err := setField(entries, index, value); err != nil {
			return fmt.Errorf("route %s: Input.%s: %s", route.Name,
				fieldName(entries, index), err.Error())
		}
	}

	if route.ResponseType != "" {
		response.SetHeader("Content-Type", route.ResponseType)
	}

	out := call(route.Handle, []any{route.Props, entries, &response})
	if len(out) == 1 && out[0] != nil {
		if failure, is := out[0].(*opinionatedagnosserver.RouteFailure); is {
			return server_run.raiseFailure(route, failure)
		}
		if handler_error, is := out[0].(error); is {
			return handler_error
		}
	}
	return nil
}

// bindValues reads every value the route declares off the request, keyed by
// the id its Input field is tagged with — a path in the type it declares,
// through the same pathValue that matched it. A required parameter the request
// does not bring, or a value that will not convert, is raised: it reports
// false, with what that failure returned.
func (server_run *run) bindValues(route *opinionatedagnosserver.Route, request serverdeps.Request) (map[string]any, bool, error) {
	values := map[string]any{fullRouteId: request.GetPath()}

	segments := splitPath(request.GetPath())
	for _, path := range route.Paths {
		text, _ := pathSlice(segments, path)
		value, _ := pathValue(path, text)
		values[path.Id] = value
	}

	for _, parameter := range route.Parameters {
		raws := parameterValues(request, parameter)

		if len(raws) == 0 {
			if parameter.Required {
				return nil, false, server_run.raise(route, opinionatedagnosserver.StatusBadRequest, parameter.Key,
					fmt.Sprintf("required parameter '%s' is missing", parameter.Key), "")
			}
			if !parameter.HasDefault {
				continue
			}
			raws = []string{parameter.Default}
		}

		value, ok := parseValue(parameter, raws)
		if !ok {
			return nil, false, server_run.raise(route, opinionatedagnosserver.StatusBadRequest, parameter.Key,
				fmt.Sprintf("parameter '%s' %s", parameter.Key, typeMessage(parameter.Type)), "")
		}
		values[parameter.Id] = value
	}

	return values, true, nil
}

// parseValue converts the raw values of one parameter to the type it
// declares, reporting false when they will not.
func parseValue(parameter opinionatedagnosserver.Parameter, raws []string) (any, bool) {
	raw := raws[0]

	switch parameter.Type {
	case opinionatedagnosserver.ParameterInteger:
		value, err := strconv.Atoi(raw)
		return value, err == nil
	case opinionatedagnosserver.ParameterIntegerArray:
		values := []int{}
		for _, one := range raws {
			value, err := strconv.Atoi(one)
			if err != nil {
				return nil, false
			}
			values = append(values, value)
		}
		return values, true
	case opinionatedagnosserver.ParameterNumber:
		value, err := strconv.ParseFloat(raw, 64)
		return value, err == nil
	case opinionatedagnosserver.ParameterBoolean:
		switch strings.ToLower(raw) {
		case "true", "1":
			return true, true
		case "false", "0":
			return false, true
		}
		return false, false
	case opinionatedagnosserver.ParameterDateTime:
		matched, err := regexp.MatchString(datetimePattern, raw)
		return raw, err == nil && matched
	case opinionatedagnosserver.ParameterStringArray:
		return raws, true
	}
	return raw, true
}

// typeMessage is how a value that will not convert is reported, in the
// server's own words rather than the conversion library's.
func typeMessage(kind opinionatedagnosserver.ParameterType) string {
	switch kind {
	case opinionatedagnosserver.ParameterInteger:
		return "is not a whole number"
	case opinionatedagnosserver.ParameterIntegerArray:
		return "holds a value that is not a whole number"
	case opinionatedagnosserver.ParameterNumber:
		return "is not a valid number"
	case opinionatedagnosserver.ParameterBoolean:
		return "must be true or false"
	case opinionatedagnosserver.ParameterDateTime:
		return "is not an RFC 3339 date-time"
	}
	return "is not valid"
}

// acceptsMediaType tells whether content_type, the Content-Type a request
// carries, is the media type route's body declares: compared without its
// parameters (";charset=utf-8") and without regard to case, never as a mere
// prefix. A request without one is accepted only when it carries no body
// either — no Content-Length above 0 and no Transfer-Encoding — so a route
// whose body is optional still takes a bare request, and a body is never read
// as a type nobody said it was.
func acceptsMediaType(route *opinionatedagnosserver.Route, request serverdeps.Request, content_type string) bool {
	if content_type == "" {
		length := request.GetHeader("Content-Length")
		return (length == "" || length == "0") && request.GetHeader("Transfer-Encoding") == ""
	}
	media_type := content_type
	if cut := strings.Index(media_type, ";"); cut >= 0 {
		media_type = media_type[:cut]
	}
	return strings.EqualFold(strings.TrimSpace(media_type), route.Body.ContentType)
}

// checkBody enforces what the body's declaration settles before a byte of it is
// read: the media type and the length the request itself declares. The body is
// the route's ReadBody's to read.
func (server_run *run) checkBody(route *opinionatedagnosserver.Route, request serverdeps.Request) (bool, error) {
	if route.Body.Type == "" || route.Body.Type == "none" {
		return true, nil
	}

	if route.Body.ContentType != "" {
		content_type := request.GetHeader("Content-Type")
		if !acceptsMediaType(route, request, content_type) {
			return false, server_run.raise(route, opinionatedagnosserver.StatusUnsupportedMediaType, "",
				fmt.Sprintf("this route accepts only a body of type %s", route.Body.ContentType), "")
		}
	}

	if raw := request.GetHeader("Content-Length"); raw != "" {
		declared, err := strconv.Atoi(raw)
		if err == nil && declared > route.Body.MaxBytes {
			return false, server_run.raise(route, opinionatedagnosserver.StatusPayloadTooLarge, "",
				fmt.Sprintf("the request body is larger than %s bytes",
					strconv.FormatInt(int64(route.Body.MaxBytes), 10)), "")
		}
	}

	return true, nil
}
