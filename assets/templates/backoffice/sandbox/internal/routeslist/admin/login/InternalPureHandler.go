package login

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficethrottle"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// InternalPureHandler answers POST /admin/login. The username field takes a
// username or an email; a match sets a session cookie bound to the client ip
// the request came from and redirects to /admin/home, anything else answers
// the login page again under a 401. Once the client ip or the login reached
// its limit of failed sign-ins, the login page is answered under a 429
// without the password being checked, until the window of
// backofficethrottle closes.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	username := entries.Body.Username

	if !backofficethrottle.LoginAllowed(sandbox, props.ClientIp, username) {
		response.SetHeader("Retry-After", backofficethrottle.RetryAfter(sandbox))
		return backofficerender.Login(sandbox, response, api.StatusTooManyRequests, "Too many failed sign-in attempts. Try again in 15 minutes.", username)
	}

	user, ok, err := backofficeauth.Authenticate(sandbox, username, entries.Body.Password)
	if err != nil {
		return err
	}
	if !ok {
		backofficethrottle.LoginFailed(sandbox, props.ClientIp, username)
		return backofficerender.Login(sandbox, response, api.StatusUnauthorized, "Invalid username or password.", username)
	}
	backofficethrottle.LoginSucceeded(sandbox, username)

	token, err := backofficeauth.IssueToken(sandbox, user, props.ClientIp)
	if err != nil {
		return err
	}

	response.AddHeader("Set-Cookie", backofficeauth.SessionCookie(sandbox, token))
	response.SetHeader("Location", "/admin/home")
	response.SetStatus(api.StatusSeeOther)
	return nil
}
