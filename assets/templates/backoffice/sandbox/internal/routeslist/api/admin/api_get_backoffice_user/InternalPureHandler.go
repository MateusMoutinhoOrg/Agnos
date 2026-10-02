package api_get_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// InternalPureHandler answers POST /api/admin/get-backoffice-user, open to
// every backoffice user like the list, with the user whose id the body names,
// or a 404 when there is none.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return routeio.Fail(sandbox, api.StatusUnauthorized, "", "no authenticated user")
	}

	user, ok := backofficeusers.Find(sandbox, int64(entries.Body.Id))
	if !ok {
		return routeio.Fail(sandbox, api.StatusNotFound, "id", "that user does not exist")
	}
	return routeio.WriteJSON(sandbox, *response, api.StatusOk, backofficeapi.UserDocument(sandbox, user))
}
