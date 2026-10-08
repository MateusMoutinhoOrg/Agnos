package backoffice_api_root_guard

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
)

// Handle runs in front of every ANY /api/admin/root/{*Rest},
// after the backoffice-api-token-auth middleware put the token's user on props.User:
// it is a middleware. A root declines, so the route after it runs; anyone
// else is refused with a 403 in JSON.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}
	if backofficeauth.Role(props.User.Role) != backofficeauth.RoleRoot {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusForbidden, "", "only root users may do this")
	}
	return nil
}
