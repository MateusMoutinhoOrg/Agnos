package revoke_backoffice_api_token_form

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapitokens"
)

// Handle answers POST /admin/revoke-backoffice-api-token/{id}:
// the token is deleted, so the api refuses it from the next request on, and
// the browser is sent to the token list with the outcome. A user revokes
// their own tokens, a root anyone's.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	notice, err := backofficeapitokens.Revoke(sandbox, *props.User, int64(input.Id))
	if err != nil {
		return err
	}
	return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficeapitokens.ListLocation(sandbox, notice))
}
