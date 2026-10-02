package list_backoffice_api_tokens

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficetokens"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// InternalPureHandler answers GET /admin/list-backoffice-api-tokens, open to
// every backoffice user, with backoffice/backoffice_api_tokens.html: the API
// tokens of the signed-in user — of every user, for a root — newest first,
// each with the control that revokes it.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return routeio.Fail(sandbox, api.StatusUnauthorized, "", "no authenticated user")
	}

	listed, err := backofficetokens.List(sandbox, *props.User)
	if err != nil {
		return err
	}
	return backofficerender.BackofficeApiTokens(sandbox, response, api.StatusOk, props.User, listed, entries.Notice, backofficerender.CreatedToken{})
}
