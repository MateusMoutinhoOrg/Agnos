package security_headers

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeguard"
)

// InternalPureHandler runs in front of every ANY /admin and /api/admin path,
// on a lower rung of the chain than the authentication middlewares: it is a
// middleware. It sets backofficeguard.SecurityHeaders and declines; setting a header
// answers nothing, so whichever route or Handle* file answers the request
// next — a page, a JSON document, a 401 — carries them.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	backofficeguard.SecurityHeaders(sandbox, response)
	return nil
}
