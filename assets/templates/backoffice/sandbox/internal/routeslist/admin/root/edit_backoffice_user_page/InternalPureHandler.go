package edit_backoffice_user_page

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// InternalPureHandler answers GET /admin/root/edit-backoffice-user/{id} with
// the form that POST /admin/root/edit-backoffice-user/{id} reads, filled with
// the user's current username, email and role. A user that does not exist sends
// the browser back to the list.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return routeio.Fail(sandbox, api.StatusUnauthorized, "", "no authenticated user")
	}

	user, ok := backofficeusers.Find(sandbox, int64(entries.Id))
	if !ok {
		return routeio.Redirect(*response, api.StatusSeeOther, backofficeusers.ListLocation(sandbox, backofficeusers.NoticeNotFound))
	}
	fields := backofficeusers.Fields{Username: user.Username, Email: user.Email, Role: user.Role}
	return backofficerender.EditBackofficeUserForm(sandbox, response, api.StatusOk, props.User, user.Id, fields, "")
}
