package backoffice_security_headers

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficehttp"
)

// Handle runs in front of every ANY /*, on a lower rung of the chain than
// the backoffice-session-auth and backoffice-api-token-auth middlewares: it is
// a middleware. On /admin and /api/admin it sets
// backofficehttp.SecurityHeaders, on every other path the
// backofficehttp.BaseSecurityHeaders a route of the application may still
// override, and declines; setting a header answers nothing, so whichever
// route or Handle* file answers the request next — a page, a JSON document, a
// 401 — carries them.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if backofficehttp.IsBackofficePath(sandbox, input.FullRoute) {
		backofficehttp.SecurityHeaders(sandbox, response)
		return nil
	}
	backofficehttp.BaseSecurityHeaders(sandbox, response)
	return nil
}
