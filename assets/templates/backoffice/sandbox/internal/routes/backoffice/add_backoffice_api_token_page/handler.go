package add_backoffice_api_token_page

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapitokens"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle answers GET /admin/add-backoffice-api-token with the
// empty form that POST /admin/add-backoffice-api-token reads, the default
// expiration selected and the client ip of the browser offered for the ips
// field.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}
	fields := backofficeapitokens.Fields{Expiration: backofficeapitokens.DefaultExpiration}
	return backofficerender.RenderAddApiTokenPage(sandbox, response, api.StatusOK, props.User, fields, props.ClientIp, "")
}
