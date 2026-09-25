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
		if len(values) == 0 || !MatchTrigger(sandbox, parameter.Trigger, values[0]) {
			return false
		}
	}

	return true
}

// MatchesPath reports whether the request path of one bound route is for it,
// whatever the method: every entry of Paths finds its slice, and the slice
// matches the entry's trigger when it declares one.
func MatchesPath(sandbox *api.Sandbox, route *api.Route) bool {
	segments := SplitPath(sandbox, routeio.RequestOf(route).GetPath())

	for _, path := range route.Paths {
		text, ok := PathSlice(sandbox, segments, path)
		if !ok {
			return false
		}
		if path.Trigger.Exist && !MatchTrigger(sandbox, path.Trigger, text) {
			return false
		}
	}

	return true
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

// MatchTrigger reports whether one text meets a trigger.
func MatchTrigger(sandbox *api.Sandbox, trigger api.Trigger, text string) bool {
	switch trigger.Type {
	case api.PrefixTrigger:
		return sandbox.Deps.Stringsdeps.HasPrefix(text, trigger.Value)
	case api.SuffixTrigger:
		return sandbox.Deps.Stringsdeps.HasSuffix(text, trigger.Value)
	case api.RegexTrigger:
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(trigger.Value, text)
		return err == nil && matched
	}
	return text == trigger.Value
}

// ParameterValues is the raw values one parameter brings, from the first of
// its Fonts that carries any: every occurrence of a query key or every
// comma-separated value of a header for a string-array, one value otherwise.
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
		if parameter.Type == api.StringArrayType {
			raws = sandbox.Deps.Stringsdeps.Split(raw, ",")
		} else {
			raws = []string{raw}
		}
	case api.QueryParam:
		if parameter.Type == api.StringArrayType {
			raws = request.GetQueryAll(parameter.Key)
		} else {
			raws = []string{request.GetQueryParam(parameter.Key)}
		}
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

// acceptsMethod reports whether a request method is one of the route's.
func acceptsMethod(route *api.Route, method string) bool {
	for _, accepted := range route.AcceptMethods {
		if accepted == method {
			return true
		}
	}
	return false
}
