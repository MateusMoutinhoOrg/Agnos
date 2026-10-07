package opinatedagnosserver

import (
	"regexp"
	"strconv"
	"strings"

	opinatedagnosserver "{{.Module}}/sandbox/deps/OpinatedAgnosServer"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// defaultIsActionable reports whether one bound route — its Request set — is
// for the request it carries: the method is one of AcceptMethods, every entry
// of Paths finds its slice and matches its trigger, and every parameter
// declaring a trigger brings a value that matches it. A parameter that is
// merely missing is not a non-match: that is the binder's to answer, with a
// 400.
//
// It is mirrored by agnos's own sandbox/internal/utils/route_match.go, which
// reads a route.yaml for `explain-route`: a change to one is a change to the
// other.
func (server_run *run) defaultIsActionable(route *opinatedagnosserver.Route) bool {
	request := route.Request

	if !acceptsMethod(route, request.GetMethod()) {
		return false
	}
	if !server_run.defaultMatchesPath(route) {
		return false
	}

	for _, parameter := range route.Parameters {
		if !parameter.Trigger.Exist {
			continue
		}
		values := parameterValues(request, parameter)
		if len(values) == 0 || !server_run.props.MatchTrigger(parameter.Trigger, values[0], false) {
			return false
		}
	}

	return true
}

// defaultMatchesPath reports whether the request path of one bound route is
// for it, whatever the method: the path has the Segments the route declares,
// and every entry of Paths finds its slice, converts to its Type and matches
// its trigger when it declares one.
func (server_run *run) defaultMatchesPath(route *opinatedagnosserver.Route) bool {
	segments := splitPath(route.Request.GetPath())

	if route.Segments > 0 && len(segments) != route.Segments {
		return false
	}

	for _, path := range route.Paths {
		text, ok := pathSlice(segments, path)
		if !ok {
			return false
		}
		if _, ok := pathValue(path, text); !ok {
			return false
		}
		if path.Trigger.Exist && !server_run.props.MatchTrigger(path.Trigger, text, true) {
			return false
		}
	}

	return true
}

// uuidPattern is what a uuid path has to read as: 8-4-4-4-12 hex digits.
const uuidPattern = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`

// pathValue converts the slice one entry of Paths read to the value its
// Entries field carries, in its Type: a path of one segment (Start equal to
// End) binds the segment itself — "42", not "/42" — and a range binds the
// slice as its trigger reads it, leading slash included. It reports false when
// the slice does not convert — which makes the route a non-match, not a bad
// request. The matcher and the binder both read a path through it, so what
// decides the match and what fills Entries never disagree.
func pathValue(path opinatedagnosserver.Path, text string) (any, bool) {
	if path.Type == opinatedagnosserver.StringPath && path.Start != path.End {
		return text, true
	}

	segment := strings.TrimPrefix(text, "/")
	switch path.Type {
	case opinatedagnosserver.IntegerPath:
		value, err := strconv.Atoi(segment)
		return value, err == nil
	case opinatedagnosserver.NumberPath:
		value, err := strconv.ParseFloat(segment, 64)
		return value, err == nil
	case opinatedagnosserver.UuidPath:
		matched, err := regexp.MatchString(uuidPattern, segment)
		return segment, err == nil && matched
	}
	return segment, true
}

// splitPath slices a raw request path into its segments, dropping the empty
// ones the leading and trailing slashes leave behind. The root path yields no
// segment at all.
func splitPath(path string) []string {
	segments := []string{}
	for _, segment := range strings.Split(path, "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}
	return segments
}

// pathSlice is the text one entry of Paths reads off the request: "/" followed
// by the segments from Start to End joined by "/", End -1 standing for the
// last one. It reports false when the request has no such slice — except the
// whole path of the root, which reads as "/".
func pathSlice(segments []string, path opinatedagnosserver.Path) (string, bool) {
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

	return "/" + strings.Join(segments[path.Start:end+1], "/"), true
}

// parameterValues is the raw values one parameter brings, from the first of
// its Fonts that carries any: every occurrence of a query key or every
// comma-separated value of a header for an array type, one value otherwise.
// It is empty when no font carries the parameter.
func parameterValues(request serverdeps.Request, parameter opinatedagnosserver.Parameter) []string {
	for _, font := range parameter.Fonts {
		values := fontValues(request, parameter, font)
		if len(values) > 0 {
			return values
		}
	}
	return []string{}
}

// fontValues is the raw values one font of the request brings for one
// parameter, empty ones dropped.
func fontValues(request serverdeps.Request, parameter opinatedagnosserver.Parameter, font opinatedagnosserver.ParameterFont) []string {
	raws := []string{}

	switch font {
	case opinatedagnosserver.HeaderParam:
		raw := request.GetHeader(parameter.Key)
		if isArrayType(parameter.Type) {
			raws = strings.Split(raw, ",")
		} else {
			raws = []string{raw}
		}
	case opinatedagnosserver.QueryParam:
		if isArrayType(parameter.Type) {
			raws = request.GetQueryAll(parameter.Key)
		} else {
			raws = []string{request.GetQueryParam(parameter.Key)}
		}
	case opinatedagnosserver.CookieParam:
		raws = []string{request.GetCookie(parameter.Key)}
	}

	values := []string{}
	for _, raw := range raws {
		value := strings.TrimSpace(raw)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}

// isArrayType reports a parameter type bound from every value the request
// brings rather than the first.
func isArrayType(kind opinatedagnosserver.ParameterType) bool {
	return kind == opinatedagnosserver.StringArrayType || kind == opinatedagnosserver.IntegerArrayType
}

// acceptsMethod reports whether a request method is one of the route's.
func acceptsMethod(route *opinatedagnosserver.Route, method string) bool {
	for _, accepted := range route.AcceptMethods {
		if accepted == method || accepted == opinatedagnosserver.AnyMethod {
			return true
		}
	}
	return false
}
