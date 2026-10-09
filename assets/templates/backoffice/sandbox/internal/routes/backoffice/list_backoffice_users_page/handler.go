package list_backoffice_users_page

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle answers GET /admin/list-backoffice-users, open to every
// backoffice user, with backoffice/backoffice_users.html: one page of the users
// the search and role filters keep, and, for a root, the controls that add,
// edit and remove them.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	listing, err := backofficeusers.List(sandbox, backofficeusers.Query{
		Search: input.Search,
		Role:   input.Role,
		Page:   input.Page,
		Limit:  input.Limit,
		Viewer: props.User,
	})
	if err != nil {
		return err
	}
	return backofficerender.RenderUsersPage(sandbox, response, props.User, listing, input.Notice)
}
