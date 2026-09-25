package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// InternalPureHandler answers {{.Methods}} {{.Trigger}}. Every value the route
// declares is already on entries, read off the request by the generic
// RequestHandler, and the response already carries the route's response-type.
//
// Setting a status is what answers the request and ends the chain. A handler
// that writes none has declined, and the next route matching this request runs
// — which is how a route becomes a middleware. Return a failure you did not
// answer yourself with routeio.Fail; nil means "done" or "not mine".
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error {
	response.SetStatus(api.StatusOk)
	response.Write([]byte("{{.Identifier}} called\n"))

	return nil
}
