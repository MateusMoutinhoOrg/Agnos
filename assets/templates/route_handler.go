package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// RouteHandler answers {{.Method}} {{.Trigger}}.
//
// Setting a status is what answers the request and ends the chain. A handler
// that writes none has declined, and the next route matching this request runs
// — which is how a route becomes a middleware. Return a failure you did not
// answer yourself with routeio.Fail; nil means "done" or "not mine".
func RouteHandler(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error {
	response.SetHeader("Content-Type", "text/plain")
	response.SetStatus(api.StatusOk)
	response.Write([]byte("{{.Identifier}} called\n"))

	return nil
}
