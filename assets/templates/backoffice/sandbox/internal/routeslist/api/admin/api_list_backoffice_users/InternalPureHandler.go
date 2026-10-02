package api_list_backoffice_users

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// InternalPureHandler answers POST /api/admin/list-backoffice-users, open to
// every backoffice user, with one page of the users the search and role
// filters keep. Every field of the body is optional, and the body itself: an
// absent page is the first one, an absent or out-of-range limit is
// backofficeusers.DefaultLimit or MaxLimit.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return routeio.Fail(sandbox, api.StatusUnauthorized, "", "no authenticated user")
	}

	listing, err := backofficeusers.List(sandbox, backofficeusers.Query{
		Search: entries.Body.Search,
		Role:   entries.Body.Role,
		Page:   entries.Body.Page,
		Limit:  entries.Body.Limit,
	})
	if err != nil {
		return err
	}
	return routeio.WriteJSON(sandbox, *response, api.StatusOk, backofficeapi.Listing(sandbox, listing))
}
