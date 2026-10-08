package list_backoffice_api_tokens_page

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapitokens"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle answers GET /admin/list-backoffice-api-tokens, open to
// every backoffice user, with backoffice/backoffice_api_tokens.html: the API
// tokens of the signed-in user — of every user, for a root — newest first,
// each with the control that revokes it.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	listed, err := backofficeapitokens.List(sandbox, *props.User)
	if err != nil {
		return err
	}
	return backofficerender.RenderApiTokensPage(sandbox, response, api.StatusOK, props.User, listed, input.Notice, backofficerender.CreatedToken{})
}
