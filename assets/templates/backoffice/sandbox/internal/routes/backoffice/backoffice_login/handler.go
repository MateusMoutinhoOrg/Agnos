package backoffice_login

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficethrottle"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle answers POST /admin/login. The username field takes a
// username or an email; a match sets a session cookie bound to the client ip
// the request came from and redirects to /admin/home, anything else answers
// the login page again under a 401. The attempt is counted before the
// password is checked, so sign-ins sent at the same time cannot all slip
// under the limit; once the client ip or the account passed its limit of
// failed sign-ins, the login page is answered under a 429 without the
// password being checked, until the window of backofficethrottle closes. A
// login or a password longer than backofficethrottle allows is refused
// unchecked.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	username := input.Body.Username
	password := input.Body.Password
	tooLong := len(username) > backofficethrottle.MaxLoginLength || len(password) > backofficethrottle.MaxPasswordLength
	if tooLong {
		username = ""
	}

	user, found, err := backofficeauth.FindUserByUsernameOrEmail(sandbox, username)
	if err != nil {
		return err
	}
	if !found {
		user.Id = 0
	}
	reservation := backofficethrottle.ReserveLogin(sandbox, props.ClientIp, backofficethrottle.AccountOf(sandbox, username, user.Id))
	if !reservation.Allowed {
		response.SetHeader("Retry-After", backofficethrottle.RetryAfter(sandbox))
		return backofficerender.RenderLoginPage(sandbox, response, api.StatusTooManyRequests, "Too many failed sign-in attempts. Try again in 15 minutes.", username)
	}
	if tooLong {
		return backofficerender.RenderLoginPage(sandbox, response, api.StatusUnauthorized, "Invalid username or password.", "")
	}

	user, ok, err := backofficeauth.CheckPassword(sandbox, user, found, password)
	if err != nil {
		return err
	}
	if !ok {
		return backofficerender.RenderLoginPage(sandbox, response, api.StatusUnauthorized, "Invalid username or password.", username)
	}
	backofficethrottle.LoginSucceeded(sandbox, reservation)

	token, err := backofficeauth.IssueSessionJWT(sandbox, user, props.ClientIp)
	if err != nil {
		return err
	}

	response.AddHeader("Set-Cookie", backofficeauth.SessionCookie(sandbox, token))
	response.SetHeader("Location", "/admin/home")
	response.SetStatus(api.StatusSeeOther)
	return nil
}
