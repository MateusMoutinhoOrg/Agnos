package api_list_backoffice_users

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// Handle answers POST /api/admin/list-backoffice-users, open to
// every backoffice user, with one page of the users the search and role
// filters keep. Every field of the body is optional, and the body itself: an
// absent page is the first one, an absent or out-of-range limit is
// backofficeusers.DefaultLimit or MaxLimit.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	listing, err := backofficeusers.List(sandbox, backofficeusers.Query{
		Search: input.Body.Search,
		Role:   input.Body.Role,
		Page:   input.Body.Page,
		Limit:  input.Body.Limit,
	})
	if err != nil {
		return err
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusOK, backofficeapi.UserListJSON(sandbox, listing))
}
