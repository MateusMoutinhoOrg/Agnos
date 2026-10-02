package remove_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// InternalPureHandler answers POST /admin/root/remove-backoffice-user/{id}: the
// user is removed, with every session of it, so its token is refused from here
// on, and the browser is sent to the list with the outcome. A root may not
// remove its own account.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return routeio.Fail(sandbox, api.StatusUnauthorized, "", "no authenticated user")
	}

	notice, err := backofficeusers.Remove(sandbox, *props.User, int64(entries.Id))
	if err != nil {
		return err
	}
	return routeio.Redirect(*response, api.StatusSeeOther, backofficeusers.ListLocation(sandbox, notice))
}
