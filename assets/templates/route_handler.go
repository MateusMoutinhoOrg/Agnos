package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
)

// Handle answers {{.Methods}} {{.Trigger}}. Every value the route
// declares is already on input, read off the request by the generic
// Run, and the response already carries the route's response-type.
//
// Answering — setting a status, or writing a byte, which sends a 200 — is what
// ends the chain. A handler that does neither has declined, and the next route
// matching this request runs — which is how a route becomes a middleware.
// What a middleware in front set on props — the request's
// routeprops.RouteProps, whose parts sandbox/internal/routeprops/ declares —
// is there to read. Refuse a request
// by returning Deps.OpinionatedAgnosServer.Fail; nil means "done" or "not mine".
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	response.SetStatus(api.StatusOK)
	response.Write([]byte("{{.Identifier}} called\n"))

	return nil
}
