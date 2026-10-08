package backoffice_same_origin

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficehttp"
)

// Handle runs in front of every ANY /admin/{*Rest}, login
// included, on a lower rung of the chain than the backoffice-session-auth middleware:
// it is a middleware.
//
// A request whose Origin header names another host than the one it was sent
// to was sent by another site's page, and is refused with a 403 before
// anything reads its session cookie: that is what a cross-site form post is.
// Anything else declines, so the route after it runs. It backs up the
// SameSite=Strict of the session cookie, which keeps another site's requests
// from carrying a session to begin with.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if backofficehttp.SameOrigin(sandbox, input.Origin, input.Host) {
		return nil
	}
	return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusForbidden, "origin", "a request sent by another site's page is refused")
}
