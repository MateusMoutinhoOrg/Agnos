package route

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeio"
)

// IsActionable reports whether one bound route — its Request set — is for the
// request it carries: the method is one of AcceptMethods, every entry of Paths
// finds its slice and matches its trigger, and every parameter declaring a
// trigger brings a value that matches it. A parameter that is merely missing is
// not a non-match: that is RequestHandler's to answer, with a 400.
func IsActionable(sandbox *api.Sandbox, route *api.Route) bool {
	request := routeio.RequestOf(route)

	if !acceptsMethod(route, request.GetMethod()) {
		return false
	}
	if !MatchesPath(sandbox, route) {
		return false
	}

	for _, parameter := range route.Parameters {
		if !parameter.Trigger.Exist {
			continue
		}
		values := ParameterValues(sandbox, request, parameter)
		if len(values) == 0 || !MatchTrigger(sandbox, parameter.Trigger, values[0], false) {
			return false
		}
	}

	return true
}

// MatchesPath reports whether the request path of one bound route is for it,
// whatever the method: the path has the Segments the route declares, and every
// entry of Paths finds its slice, converts to its Type and matches its trigger
// when it declares one.
func MatchesPath(sandbox *api.Sandbox, route *api.Route) bool {
	segments := SplitPath(sandbox, routeio.RequestOf(route).GetPath())

	if route.Segments > 0 && len(segments) != route.Segments {
		return false
	}

	for _, path := range route.Paths {
		text, ok := PathSlice(sandbox, segments, path)
		if !ok {
			return false
		}
		if _, ok := PathValue(sandbox, path, text); !ok {
			return false
		}
		if path.Trigger.Exist && !MatchTrigger(sandbox, path.Trigger, text, true) {
			return false
		}
	}

	return true
}

// uuidPattern is what a uuid path has to read as: 8-4-4-4-12 hex digits.
const uuidPattern = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`

// PathValue converts the slice one entry of Paths read to the value its
// Entries field carries, in its Type: a path of one segment (Start equal to
// End) binds the segment itself — "42", not "/42" — and a range binds the
// slice as its trigger reads it, leading slash included. It reports false when
// the slice does not convert — which makes the route a non-match, not a bad
// request. MatchesPath and RequestHandler both read a path through it, so what
// decides the match and what fills Entries never disagree.
func PathValue(sandbox *api.Sandbox, path api.Path, text string) (any, bool) {
	if path.Type == api.StringPath && path.Start != path.End {
		return text, true
	}

	segment := sandbox.Deps.Stringsdeps.TrimPrefix(text, "/")
	switch path.Type {
	case api.IntegerPath:
		value, err := sandbox.Deps.Stringsdeps.Atoi(segment)
		return value, err == nil
	case api.NumberPath:
		value, err := sandbox.Deps.Stringsdeps.ParseFloat(segment, 64)
		return value, err == nil
	case api.UuidPath:
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(uuidPattern, segment)
		return segment, err == nil && matched
	}
	return segment, true
}

// SplitPath slices a raw request path into its segments, dropping the empty
// ones the leading and trailing slashes leave behind. The root path yields no
// segment at all.
func SplitPath(sandbox *api.Sandbox, path string) []string {
	segments := []string{}
	for _, segment := range sandbox.Deps.Stringsdeps.Split(path, "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}
	return segments
}

// PathSlice is the text one entry of Paths reads off the request: "/" followed
// by the segments from Start to End joined by "/", End -1 standing for the
// last one. It reports false when the request has no such slice — except the
// whole path of the root, which reads as "/".
func PathSlice(sandbox *api.Sandbox, segments []string, path api.Path) (string, bool) {
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

// MatchTrigger reports whether one text meets a trigger: the comparison its
// Type names, run without regard to case when it declares IgnoreCase, and
// inverted when it declares Negate. Segmented is true for a path slice, the
// one text a prefix reads segment by segment.
func MatchTrigger(sandbox *api.Sandbox, trigger api.Trigger, text string, segmented bool) bool {
	return compareTrigger(sandbox, trigger, text, segmented) != trigger.Negate
}

// compareTrigger is the comparison of MatchTrigger before Negate. On a path
// slice a prefix holds on a segment boundary alone — "/admin" is "/admin" or
// "/admin/…", never "/administrator" — which is what tells it from a
// text-prefix. A parameter value has no segments, so there the two are one.
func compareTrigger(sandbox *api.Sandbox, trigger api.Trigger, text string, segmented bool) bool {
	value := trigger.Value
	if trigger.IgnoreCase && trigger.Type != api.RegexTrigger {
		value = sandbox.Deps.Stringsdeps.ToLower(value)
		text = sandbox.Deps.Stringsdeps.ToLower(text)
	}

	switch trigger.Type {
	case api.PrefixTrigger:
		if !segmented {
			return sandbox.Deps.Stringsdeps.HasPrefix(text, value)
		}
		value = sandbox.Deps.Stringsdeps.TrimSuffix(value, "/")
		return value == "" || text == value || sandbox.Deps.Stringsdeps.HasPrefix(text, value+"/")
	case api.TextPrefixTrigger:
		return sandbox.Deps.Stringsdeps.HasPrefix(text, value)
	case api.SuffixTrigger:
		return sandbox.Deps.Stringsdeps.HasSuffix(text, value)
	case api.RegexTrigger:
		if trigger.IgnoreCase {
			value = "(?i)" + value
		}
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(value, text)
		return err == nil && matched
	}
	return text == value
}

// ParameterValues is the raw values one parameter brings, from the first of
// its Fonts that carries any: every occurrence of a query key or every
// comma-separated value of a header for an array type, one value otherwise.
// It is empty when no font carries the parameter.
func ParameterValues(sandbox *api.Sandbox, request serverdeps.Request, parameter api.Parameter) []string {
	for _, font := range parameter.Fonts {
		values := fontValues(sandbox, request, parameter, font)
		if len(values) > 0 {
			return values
		}
	}
	return []string{}
}

// fontValues is the raw values one font of the request brings for one
// parameter, empty ones dropped.
func fontValues(sandbox *api.Sandbox, request serverdeps.Request, parameter api.Parameter, font api.ParameterFont) []string {
	raws := []string{}

	switch font {
	case api.HeaderParam:
		raw := request.GetHeader(parameter.Key)
		if isArrayType(parameter.Type) {
			raws = sandbox.Deps.Stringsdeps.Split(raw, ",")
		} else {
			raws = []string{raw}
		}
	case api.QueryParam:
		if isArrayType(parameter.Type) {
			raws = request.GetQueryAll(parameter.Key)
		} else {
			raws = []string{request.GetQueryParam(parameter.Key)}
		}
	case api.CookieParam:
		raws = []string{request.GetCookie(parameter.Key)}
	}

	values := []string{}
	for _, raw := range raws {
		value := sandbox.Deps.Stringsdeps.TrimSpace(raw)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}

// isArrayType reports a parameter type bound from every value the request
// brings rather than the first.
func isArrayType(kind api.ParameterType) bool {
	return kind == api.StringArrayType || kind == api.IntegerArrayType
}

// acceptsMethod reports whether a request method is one of the route's.
func acceptsMethod(route *api.Route, method string) bool {
	for _, accepted := range route.AcceptMethods {
		if accepted == method || accepted == api.AnyMethod {
			return true
		}
	}
	return false
}
