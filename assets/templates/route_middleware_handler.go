package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
)

// Handle runs in front of every {{.Methods}} {{.Trigger}} on a
// lower rung of the chain: it is a middleware. Returning nil without answering
// hands the request to the next route; answering — a status, or a byte —
// ends the chain here.
//
// Refuse a request by returning Deps.OpinionatedAgnosServer.Fail, which is answered through the
// project's own handler for that status:
//
//	return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "authorization", "invalid token")
//
// Hand what you learned to the routes after it through props, the request's
// routeprops.RouteProps — declare the field in sandbox/internal/routeprops/project.go:
//
//	props.User = user
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	return nil
}
