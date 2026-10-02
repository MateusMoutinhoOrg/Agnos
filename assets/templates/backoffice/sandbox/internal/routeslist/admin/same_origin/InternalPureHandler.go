package same_origin

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeguard"
)

// InternalPureHandler runs in front of every ANY /admin/{*Rest}, login
// included, on a lower rung of the chain than the authentication middleware:
// it is a middleware.
//
// A request whose Origin header names another host than the one it was sent
// to was sent by another site's page, and is refused with a 403 before
// anything reads its session cookie: that is what a cross-site form post is.
// Anything else declines, so the route after it runs. It backs up the
// SameSite=Strict of the session cookie, which keeps another site's requests
// from carrying a session to begin with.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if backofficeguard.SameOrigin(sandbox, entries.Origin, entries.Host) {
		return nil
	}
	return routeio.Fail(sandbox, api.StatusForbidden, "origin", "a request sent by another site's page is refused")
}
