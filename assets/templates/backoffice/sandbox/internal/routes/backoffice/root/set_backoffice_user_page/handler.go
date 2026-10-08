package set_backoffice_user_page

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle answers GET /admin/root/set-backoffice-user/{id} with
// the form that POST /admin/root/set-backoffice-user/{id} reads, filled with
// the user's current username, email and role. A user that does not exist sends
// the browser back to the list.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	user, ok := backofficeusers.Find(sandbox, int64(input.Id))
	if !ok {
		return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficeusers.ListLocation(sandbox, backofficeusers.NoticeNotFound))
	}
	fields := backofficeusers.Fields{Username: user.Username, Email: user.Email, Role: user.Role}
	return backofficerender.RenderSetUserPage(sandbox, response, api.StatusOK, props.User, user.Id, fields, "")
}
