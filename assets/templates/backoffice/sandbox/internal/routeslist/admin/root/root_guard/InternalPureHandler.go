package root_guard

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// InternalPureHandler runs in front of every ANY /admin/root/{*Rest}, after the
// authentication middleware put the signed-in user on props.User: it is a
// middleware. A root declines, so the route after it runs; anyone else is
// answered the forbidden page under a 403.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return routeio.Fail(sandbox, api.StatusUnauthorized, "", "no authenticated user")
	}
	if backofficeauth.Role(props.User.Role) != backofficeauth.RoleRoot {
		return backofficerender.Forbidden(sandbox, response, props.User)
	}
	return nil
}
