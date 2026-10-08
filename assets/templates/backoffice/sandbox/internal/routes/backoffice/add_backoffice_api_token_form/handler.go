package add_backoffice_api_token_form

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapitokens"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle answers POST /admin/add-backoffice-api-token. A token
// the form describes well is created for the signed-in user and the token list
// is answered with it shown in full above the list — the one time it is ever
// shown, which is why this answers the page instead of redirecting: the token
// never travels in a url. Anything else answers the form again, filled with
// what was sent, under a 400 with the reason above it.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	fields := backofficeapitokens.Fields{
		Name:       input.Body.Name,
		Expiration: input.Body.Expiration,
		Date:       input.Body.Date,
		Ips:        input.Body.Ips,
	}
	token, item, message, err := backofficeapitokens.Add(sandbox, *props.User, fields)
	if err != nil {
		return err
	}
	if message != "" {
		return backofficerender.RenderAddApiTokenPage(sandbox, response, api.StatusBadRequest, props.User, fields, props.ClientIp, message)
	}

	listed, err := backofficeapitokens.List(sandbox, *props.User)
	if err != nil {
		return err
	}
	return backofficerender.RenderApiTokensPage(sandbox, response, api.StatusCreated, props.User, listed, "", backofficerender.CreatedToken{Name: item.Name, Token: token})
}
