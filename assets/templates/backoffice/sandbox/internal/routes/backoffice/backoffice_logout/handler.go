package backoffice_logout

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
)

// Handle answers POST /admin/logout for the session the
// backoffice-session-auth middleware put on props.Session: the session is closed, so its
// token is refused from here on, the session cookie is cleared, and the
// browser is sent back to /admin/home, which answers the login page. Every
// other session of the user, on this device or another, stays open.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil || props.Session == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	err := backofficeauth.Logout(sandbox, *props.User, *props.Session)
	if err != nil {
		return err
	}

	response.AddHeader("Set-Cookie", backofficeauth.ClearedCookie(sandbox))
	response.SetHeader("Location", "/admin/home")
	response.SetStatus(api.StatusSeeOther)
	return nil
}
