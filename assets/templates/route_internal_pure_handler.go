package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// InternalPureHandler answers {{.Methods}} {{.Trigger}}. Every value the route
// declares is already on entries, read off the request by the generic
// RequestHandler, and the response already carries the route's response-type.
//
// Answering — setting a status, or writing a byte, which sends a 200 — is what
// ends the chain. A handler that does neither has declined, and the next route
// matching this request runs — which is how a route becomes a middleware.
// What a middleware in front set on props — the request's api.RouteProps,
// declared in sandbox/api/routeprops.go — is there to read. Refuse a request
// by returning routeio.Fail; nil means "done" or "not mine".
func InternalPureHandler(sandbox *api.Sandbox, props *api.RouteProps, entries *Entries, response *serverdeps.Response) error {
	response.SetStatus(api.StatusOk)
	response.Write([]byte("{{.Identifier}} called\n"))

	return nil
}
