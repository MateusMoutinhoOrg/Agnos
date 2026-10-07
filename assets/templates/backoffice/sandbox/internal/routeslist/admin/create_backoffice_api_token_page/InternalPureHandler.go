package create_backoffice_api_token_page

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficetokens"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// InternalPureHandler answers GET /admin/create-backoffice-api-token with the
// empty form that POST /admin/create-backoffice-api-token reads, the default
// expiration selected and the client ip of the browser offered for the ips
// field.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}
	fields := backofficetokens.Fields{Expiration: backofficetokens.DefaultExpiration}
	return backofficerender.CreateBackofficeApiTokenForm(sandbox, response, api.StatusOk, props.User, fields, props.ClientIp, "")
}
