package backoffice_security_headers

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficehttp"
)

// Handle runs in front of every ANY /admin and /api/admin path,
// on a lower rung of the chain than the backoffice-session-auth and backoffice-api-token-auth middlewares: it is a
// middleware. It sets backofficehttp.SecurityHeaders and declines; setting a header
// answers nothing, so whichever route or Handle* file answers the request
// next — a page, a JSON document, a 401 — carries them.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	backofficehttp.SecurityHeaders(sandbox, response)
	return nil
}
