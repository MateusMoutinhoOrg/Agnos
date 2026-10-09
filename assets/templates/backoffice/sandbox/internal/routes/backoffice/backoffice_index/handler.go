package backoffice_index

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
)

// Handle answers GET /admin, the backoffice's own root, with a 303 to
// /admin/home for the user the backoffice-session-auth middleware put on
// props.User. A request without a session never reaches it: the middleware
// answers the login page first.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}
	response.SetHeader("Location", "/admin/home")
	response.SetStatus(api.StatusSeeOther)
	return nil
}
