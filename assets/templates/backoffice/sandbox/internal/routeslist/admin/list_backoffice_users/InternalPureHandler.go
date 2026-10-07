package list_backoffice_users

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// InternalPureHandler answers GET /admin/list-backoffice-users, open to every
// backoffice user, with backoffice/backoffice_users.html: one page of the users
// the search and role filters keep, and, for a root, the controls that add,
// edit and remove them.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	listing, err := backofficeusers.List(sandbox, backofficeusers.Query{
		Search: entries.Search,
		Role:   entries.Role,
		Page:   entries.Page,
		Limit:  entries.Limit,
	})
	if err != nil {
		return err
	}
	return backofficerender.BackofficeUsers(sandbox, response, props.User, listing, entries.Notice)
}
