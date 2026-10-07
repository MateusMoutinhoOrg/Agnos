package api_get_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// InternalPureHandler answers POST /api/admin/get-backoffice-user, open to
// every backoffice user like the list, with the user whose id the body names,
// or a 404 when there is none.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	user, ok := backofficeusers.Find(sandbox, int64(entries.Body.Id))
	if !ok {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusNotFound, "id", "that user does not exist")
	}
	return sandbox.Deps.OpinatedAgnosServer.WriteJSON(sandbox.Deps.Serializables, *response, api.StatusOk, backofficeapi.UserDocument(sandbox, user))
}
