package create_backoffice_api_token

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficetokens"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// InternalPureHandler answers POST /admin/create-backoffice-api-token. A token
// the form describes well is created for the signed-in user and the token list
// is answered with it shown in full above the list — the one time it is ever
// shown, which is why this answers the page instead of redirecting: the token
// never travels in a url. Anything else answers the form again, filled with
// what was sent, under a 400 with the reason above it.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	fields := backofficetokens.Fields{
		Name:       entries.Body.Name,
		Expiration: entries.Body.Expiration,
		Date:       entries.Body.Date,
		Ips:        entries.Body.Ips,
	}
	token, item, message, err := backofficetokens.Create(sandbox, *props.User, fields)
	if err != nil {
		return err
	}
	if message != "" {
		return backofficerender.CreateBackofficeApiTokenForm(sandbox, response, api.StatusBadRequest, props.User, fields, props.ClientIp, message)
	}

	listed, err := backofficetokens.List(sandbox, *props.User)
	if err != nil {
		return err
	}
	return backofficerender.BackofficeApiTokens(sandbox, response, api.StatusCreated, props.User, listed, "", backofficerender.CreatedToken{Name: item.Name, Token: token})
}
