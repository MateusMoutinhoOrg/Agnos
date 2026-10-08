package backoffice_session_auth

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle runs in front of every ANY /admin/{*Rest} except
// /admin/login, on a lower rung of the chain: it is a middleware.
//
// A valid session cookie, issued to the client ip this request came from —
// props.ClientIp, which the backoffice-client-ip middleware worked out — for a
// session no logout has closed, puts its user on props.User and its session on
// props.Session and declines, so the route after it runs. Anything else answers the login page under a 401, and
// clears a cookie that no longer holds a valid session.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	user, session, ok := backofficeauth.ResolveSession(sandbox, input.AdminToken, props.ClientIp)
	if ok {
		props.User = &user
		props.Session = &session
		return nil
	}

	message := ""
	if input.AdminToken != "" {
		response.AddHeader("Set-Cookie", backofficeauth.ClearedCookie(sandbox))
		message = "Your session has expired. Please sign in again."
	}
	return backofficerender.RenderLoginPage(sandbox, response, api.StatusUnauthorized, message, "")
}
