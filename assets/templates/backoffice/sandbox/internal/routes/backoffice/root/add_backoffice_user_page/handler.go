package add_backoffice_user_page

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// InternalPureHandler answers GET /admin/root/add-backoffice-user with the
// empty form that POST /admin/root/add-backoffice-user reads, the viewer role
// selected.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}
	fields := backofficeusers.Fields{Role: int64(backofficeauth.RoleViewer)}
	return backofficerender.AddBackofficeUserForm(sandbox, response, api.StatusOk, props.User, fields, "")
}
