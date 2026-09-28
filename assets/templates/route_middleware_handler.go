package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// InternalPureHandler runs in front of every {{.Methods}} {{.Trigger}} on a
// lower rung of the chain: it is a middleware. Returning nil without answering
// hands the request to the next route; answering — a status, or a byte —
// ends the chain here.
//
// Refuse a request by returning routeio.Fail, which is answered through the
// project's own handler for that status:
//
//	return routeio.Fail(sandbox, api.StatusUnauthorized, "authorization", "invalid token")
//
// Hand what you learned to the routes after it through props, the request's
// api.RouteProps — declare the field in sandbox/api/routeprops.go:
//
//	props.User = user
func InternalPureHandler(sandbox *api.Sandbox, props *api.RouteProps, entries *Entries, response *serverdeps.Response) error {
	return nil
}
