package explain_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ExplainRouteInternal walks the chain the way the generated dispatch does and
// says, route by route, whether it runs for the request and why not when it
// does not. What a handler does once it runs is its own code, so the answer is
// read as "runs": the first route that answers ends the chain, and a route
// that answers nothing hands it on. The last line is what the request ends on
// when no route answers — the 405 or the 404 the dispatch raises.
func ExplainRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.ExplainRouteProps) ([]string, error) {
	if !io.IsDir(utils.RoutesDir) {
		return nil, sandbox.Deps.Std.Errorf("the project has no server layer: run server-init first")
	}

	request, err := requestOf(sandbox, props)
	if err != nil {
		return nil, err
	}

	chain, err := utils.LoadRouteChain(sandbox, io)
	if err != nil {
		return nil, err
	}

	lines := []string{sandbox.Deps.Std.Sprintf("%s %s", request.Method, props.RequestPath)}
	ran := false
	method_mismatch := false
	answered := ""
	phase := routeconf.PhaseBefore

	for _, entry := range chain {
		if entry.Conf.Phase != phase {
			phase = entry.Conf.Phase
			lines = append(lines, "", "after the answer")
		}

		match := utils.MatchRouteRequest(sandbox, entry.Conf, request)
		name := utils.RouteIdentifier(sandbox, entry.Name)
		head := sandbox.Deps.Std.Sprintf("  %-4d %-24s", entry.Conf.Priority, name)

		if !match.Runs {
			if match.MethodMismatch && phase == routeconf.PhaseBefore {
				method_mismatch = true
			}
			lines = append(lines, head+" skipped: "+match.Reason)
			continue
		}

		if phase == routeconf.PhaseBefore && answered != "" {
			lines = append(lines, head+" not reached: "+answered+" answered first")
			continue
		}
		if phase == routeconf.PhaseBefore && !isAny(entry.Conf) {
			ran = true
		}
		if match.Failure != "" && phase == routeconf.PhaseBefore {
			lines = append(lines, head+" answers "+match.Failure+", before its handler runs")
			answered = name
			continue
		}
		lines = append(lines, head+" runs")
	}

	lines = append(lines, "")
	switch {
	case answered != "":
		lines = append(lines, "the request ends on "+answered)
	case ran:
		lines = append(lines, "the first route that answers ends the chain; if none does: 404")
	case request.Method == "HEAD":
		lines = append(lines, "nothing declares HEAD for it: the dispatch runs the chain again as a GET, and answers 404 when nothing answers that")
	case method_mismatch:
		lines = append(lines, "if no route answers: 405, a route takes the path under another method")
	default:
		lines = append(lines, "if no route answers: 404")
	}

	return lines, nil
}

// isAny reports a route declared for every method, which the dispatch does
// not count when it tells a 405 from a 404.
func isAny(conf *routeconf.RouteConf) bool {
	return len(conf.Methods) == 1 && conf.Methods[0] == routeconf.AnyMethod
}

// requestOf reads the request typed on the command line: the method (GET when
// none), the path with its query string, and "key=value" headers and cookies.
func requestOf(sandbox *api.Sandbox, props api.ExplainRouteProps) (utils.RouteRequest, error) {
	method := sandbox.Deps.Stringsdeps.ToUpper(sandbox.Deps.Stringsdeps.TrimSpace(props.Method))
	if method == "" {
		method = routeconf.DefaultMethod
	}
	if method != routeconf.AnyMethod {
		if _, err := utils.RouteMethod(sandbox, method); err != nil {
			return utils.RouteRequest{}, err
		}
	}

	raw := sandbox.Deps.Stringsdeps.TrimSpace(props.RequestPath)
	if !sandbox.Deps.Stringsdeps.HasPrefix(raw, "/") {
		return utils.RouteRequest{}, sandbox.Deps.Std.Errorf("the request path %q has to start with /", props.RequestPath)
	}

	request := utils.RouteRequest{
		Method:  method,
		Path:    raw,
		Query:   map[string][]string{},
		Headers: map[string]string{},
		Cookies: map[string]string{},
	}

	if parts := splitOnce(sandbox, raw, "?"); len(parts) == 2 {
		request.Path = parts[0]
		for _, pair := range sandbox.Deps.Stringsdeps.Split(parts[1], "&") {
			if pair == "" {
				continue
			}
			key, value := splitPair(sandbox, pair)
			request.Query[key] = append(request.Query[key], value)
		}
	}

	for _, header := range props.Headers {
		key, value, err := typedPair(sandbox, "--header", header)
		if err != nil {
			return request, err
		}
		request.Headers[sandbox.Deps.Stringsdeps.ToLower(key)] = value
	}
	for _, cookie := range props.Cookies {
		key, value, err := typedPair(sandbox, "--cookie", cookie)
		if err != nil {
			return request, err
		}
		request.Cookies[key] = value
	}

	return request, nil
}

// splitPair reads one key=value, a bare key reading as an empty value.
func splitPair(sandbox *api.Sandbox, pair string) (string, string) {
	parts := splitOnce(sandbox, pair, "=")
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

// typedPair reads one key=value typed after a flag, refusing one with no key.
func typedPair(sandbox *api.Sandbox, flag string, pair string) (string, string, error) {
	key, value := splitPair(sandbox, pair)
	key = sandbox.Deps.Stringsdeps.TrimSpace(key)
	if key == "" || !sandbox.Deps.Stringsdeps.Contains(pair, "=") {
		return "", "", sandbox.Deps.Std.Errorf("%s %q is not key=value", flag, pair)
	}
	return key, sandbox.Deps.Stringsdeps.TrimSpace(value), nil
}

// splitOnce splits text at the first sep: one part when there is none, two
// otherwise.
func splitOnce(sandbox *api.Sandbox, text string, sep string) []string {
	parts := sandbox.Deps.Stringsdeps.Split(text, sep)
	if len(parts) < 2 {
		return parts
	}
	return []string{parts[0], sandbox.Deps.Stringsdeps.Join(parts[1:], sep)}
}
