package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)
{{- if .After}}

// InternalPureHandler runs after every {{.Methods}} {{.Trigger}} has been
// answered, whatever answered it: the route is in the `after` phase. The
// response is frozen — a status, a header or a byte written here is dropped —
// so this is the place for what reads the answer rather than makes it: a log
// line, a metric. routeio.AnsweredStatus(route) is the status it went out
// with, and route.Locals holds what the chain stored.
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error {
	return nil
}
{{- else}}

// InternalPureHandler runs in front of every {{.Methods}} {{.Trigger}} on a
// lower rung of the chain: it is a middleware. Returning nil without answering
// hands the request to the next route; answering — a status, or a byte —
// ends the chain here.
//
// Refuse a request with routeio.Fail, which answers it through the project's
// own handler for that status:
//
//	return routeio.Fail(sandbox, route, api.StatusUnauthorized, "authorization", "invalid token")
//
// Hand what you learned to the routes after it through route.Locals:
//
//	routeio.SetLocal(route, "user", user)
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error {
	return nil
}
{{- end}}
