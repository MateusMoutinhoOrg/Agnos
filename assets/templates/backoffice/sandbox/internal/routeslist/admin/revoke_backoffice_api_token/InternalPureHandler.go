package revoke_backoffice_api_token

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficetokens"
)

// InternalPureHandler answers POST /admin/revoke-backoffice-api-token/{id}:
// the token is deleted, so the api refuses it from the next request on, and
// the browser is sent to the token list with the outcome. A user revokes
// their own tokens, a root anyone's.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	notice, err := backofficetokens.Revoke(sandbox, *props.User, int64(entries.Id))
	if err != nil {
		return err
	}
	return sandbox.Deps.OpinatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficetokens.ListLocation(sandbox, notice))
}
